package tests

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/yanmhlv/pcloud"
)

func getClient(t *testing.T) (*pcloud.Client, context.Context) {
	t.Helper()

	username := os.Getenv("PCLOUD_USERNAME")
	password := os.Getenv("PCLOUD_PASSWORD")
	if username == "" || password == "" {
		t.Skip("PCLOUD_USERNAME and PCLOUD_PASSWORD env vars required")
	}

	baseURL := os.Getenv("PCLOUD_BASE_URL")
	c := pcloud.NewClient(baseURL)
	ctx := context.Background()

	if err := c.Login(ctx, username, password); err != nil {
		t.Fatalf("login failed: %v", err)
	}

	return c, ctx
}

func TestAuth(t *testing.T) {
	username := os.Getenv("PCLOUD_USERNAME")
	password := os.Getenv("PCLOUD_PASSWORD")
	if username == "" || password == "" {
		t.Skip("PCLOUD_USERNAME and PCLOUD_PASSWORD env vars required")
	}

	baseURL := os.Getenv("PCLOUD_BASE_URL")
	c := pcloud.NewClient(baseURL)
	ctx := context.Background()

	t.Run("Login", func(t *testing.T) {
		if err := c.Login(ctx, username, password); err != nil {
			t.Fatalf("login failed: %v", err)
		}
	})

	t.Run("UserInfo", func(t *testing.T) {
		info, err := c.UserInfo(ctx)
		if err != nil {
			t.Fatalf("userinfo failed: %v", err)
		}
		if info.Email == "" {
			t.Fatal("email is empty")
		}
		if info.Quota == 0 {
			t.Fatal("quota is zero")
		}
	})

	t.Run("Logout", func(t *testing.T) {
		if err := c.Logout(ctx); err != nil {
			t.Fatalf("logout failed: %v", err)
		}
	})
}

func TestFolders(t *testing.T) {
	c, ctx := getClient(t)
	defer c.Logout(ctx)

	testFolder := "pcloud_test_" + time.Now().Format("20060102150405")
	var folderID uint64

	t.Run("CreateFolder", func(t *testing.T) {
		folder, err := c.CreateFolder(ctx, 0, testFolder)
		if err != nil {
			t.Fatalf("create folder failed: %v", err)
		}
		if folder.Name != testFolder {
			t.Fatalf("expected name %s, got %s", testFolder, folder.Name)
		}
		folderID = folder.FolderID
	})

	t.Run("ListFolder", func(t *testing.T) {
		folder, err := c.ListFolder(ctx, 0, nil)
		if err != nil {
			t.Fatalf("list folder failed: %v", err)
		}
		if !folder.IsFolder {
			t.Fatal("root should be a folder")
		}
	})

	t.Run("CreateFolderIfNotExists", func(t *testing.T) {
		folder, err := c.CreateFolderIfNotExists(ctx, folderID, "subfolder")
		if err != nil {
			t.Fatalf("create subfolder failed: %v", err)
		}
		if folder.Name != "subfolder" {
			t.Fatalf("expected name subfolder, got %s", folder.Name)
		}

		folder2, err := c.CreateFolderIfNotExists(ctx, folderID, "subfolder")
		if err != nil {
			t.Fatalf("create existing subfolder failed: %v", err)
		}
		if folder2.FolderID != folder.FolderID {
			t.Fatal("should return same folder")
		}
	})

	t.Run("RenameFolder", func(t *testing.T) {
		newName := testFolder + "_renamed"
		folder, err := c.RenameFolder(ctx, folderID, newName)
		if err != nil {
			t.Fatalf("rename folder failed: %v", err)
		}
		if folder.Name != newName {
			t.Fatalf("expected name %s, got %s", newName, folder.Name)
		}

		folder, _ = c.RenameFolder(ctx, folderID, testFolder)
		if folder.Name != testFolder {
			t.Fatal("rename back failed")
		}
	})

	t.Run("CopyFolder", func(t *testing.T) {
		copyName := testFolder + "_copy"
		copyFolder, err := c.CreateFolder(ctx, 0, copyName)
		if err != nil {
			t.Fatalf("create copy target failed: %v", err)
		}
		defer c.DeleteFolderRecursive(ctx, copyFolder.FolderID)

		subFolder, _ := c.CreateFolder(ctx, folderID, "to_copy")
		copied, err := c.CopyFolder(ctx, subFolder.FolderID, copyFolder.FolderID)
		if err != nil {
			t.Fatalf("copy folder failed: %v", err)
		}
		if copied.Name != "to_copy" {
			t.Fatal("copied folder has wrong name")
		}
	})

	t.Run("DeleteFolderRecursive", func(t *testing.T) {
		if err := c.DeleteFolderRecursive(ctx, folderID); err != nil {
			t.Fatalf("delete folder recursive failed: %v", err)
		}
	})
}

func TestFiles(t *testing.T) {
	c, ctx := getClient(t)
	defer c.Logout(ctx)

	testFolder := "pcloud_file_test_" + time.Now().Format("20060102150405")
	folder, err := c.CreateFolder(ctx, 0, testFolder)
	if err != nil {
		t.Fatalf("create test folder failed: %v", err)
	}
	defer c.DeleteFolderRecursive(ctx, folder.FolderID)

	testContent := []byte("hello pcloud test content")
	var fileID uint64

	t.Run("Upload", func(t *testing.T) {
		meta, err := c.Upload(ctx, folder.FolderID, "test.txt", bytes.NewReader(testContent), nil)
		if err != nil {
			t.Fatalf("upload failed: %v", err)
		}
		if meta.Name != "test.txt" {
			t.Fatalf("expected name test.txt, got %s", meta.Name)
		}
		if meta.Size != uint64(len(testContent)) {
			t.Fatalf("expected size %d, got %d", len(testContent), meta.Size)
		}
		fileID = meta.FileID
	})

	t.Run("Stat", func(t *testing.T) {
		meta, err := c.Stat(ctx, fileID)
		if err != nil {
			t.Fatalf("stat failed: %v", err)
		}
		if meta.Name != "test.txt" {
			t.Fatalf("expected name test.txt, got %s", meta.Name)
		}
	})

	t.Run("Download", func(t *testing.T) {
		body, err := c.Download(ctx, fileID)
		if err != nil {
			t.Fatalf("download failed: %v", err)
		}
		defer body.Close()

		content, err := io.ReadAll(body)
		if err != nil {
			t.Fatalf("read body failed: %v", err)
		}
		if !bytes.Equal(content, testContent) {
			t.Fatalf("content mismatch: got %s", content)
		}
	})

	t.Run("RenameFile", func(t *testing.T) {
		meta, err := c.RenameFile(ctx, fileID, "renamed.txt")
		if err != nil {
			t.Fatalf("rename file failed: %v", err)
		}
		if meta.Name != "renamed.txt" {
			t.Fatalf("expected name renamed.txt, got %s", meta.Name)
		}
	})

	t.Run("CopyFile", func(t *testing.T) {
		copyTarget, err := c.CreateFolder(ctx, folder.FolderID, "copy_target")
		if err != nil {
			t.Fatalf("create copy target folder failed: %v", err)
		}
		meta, err := c.CopyFile(ctx, fileID, copyTarget.FolderID)
		if err != nil {
			t.Fatalf("copy file failed: %v", err)
		}
		if meta.FileID == fileID {
			t.Fatal("copied file should have different id")
		}
		c.DeleteFile(ctx, meta.FileID)
	})

	t.Run("DeleteFile", func(t *testing.T) {
		if err := c.DeleteFile(ctx, fileID); err != nil {
			t.Fatalf("delete file failed: %v", err)
		}

		_, err := c.Stat(ctx, fileID)
		if err == nil {
			t.Fatal("stat should fail for deleted file")
		}
	})
}

func TestRevisions(t *testing.T) {
	c, ctx := getClient(t)
	defer c.Logout(ctx)

	testFolder := "pcloud_rev_test_" + time.Now().Format("20060102150405")
	folder, err := c.CreateFolder(ctx, 0, testFolder)
	if err != nil {
		t.Fatalf("create test folder failed: %v", err)
	}
	defer c.DeleteFolderRecursive(ctx, folder.FolderID)

	content1 := []byte("version 1")
	content2 := []byte("version 2")

	meta, err := c.Upload(ctx, folder.FolderID, "revtest.txt", bytes.NewReader(content1), nil)
	if err != nil {
		t.Fatalf("upload v1 failed: %v", err)
	}
	fileID := meta.FileID

	_, err = c.Upload(ctx, folder.FolderID, "revtest.txt", bytes.NewReader(content2), nil)
	if err != nil {
		t.Fatalf("upload v2 failed: %v", err)
	}

	t.Run("ListRevisions", func(t *testing.T) {
		revs, err := c.ListRevisions(ctx, fileID)
		if err != nil {
			t.Fatalf("list revisions failed: %v", err)
		}
		if len(revs) == 0 {
			t.Fatal("expected at least 1 revision")
		}
	})

	t.Run("RevertRevision", func(t *testing.T) {
		revs, _ := c.ListRevisions(ctx, fileID)
		if len(revs) < 2 {
			t.Skip("not enough revisions")
		}

		oldRev := revs[len(revs)-1]
		_, err := c.RevertRevision(ctx, fileID, oldRev.RevisionID)
		if err != nil {
			t.Fatalf("revert revision failed: %v", err)
		}

		body, _ := c.Download(ctx, fileID)
		defer body.Close()
		content, _ := io.ReadAll(body)
		if !bytes.Equal(content, content1) {
			t.Fatalf("content should be reverted to v1, got %s", content)
		}
	})
}

func TestStreaming(t *testing.T) {
	c, ctx := getClient(t)
	defer c.Logout(ctx)

	testFolder := "pcloud_stream_test_" + time.Now().Format("20060102150405")
	folder, err := c.CreateFolder(ctx, 0, testFolder)
	if err != nil {
		t.Fatalf("create test folder failed: %v", err)
	}
	defer c.DeleteFolderRecursive(ctx, folder.FolderID)

	content := []byte("stream test content")
	meta, err := c.Upload(ctx, folder.FolderID, "stream.txt", bytes.NewReader(content), nil)
	if err != nil {
		t.Fatalf("upload failed: %v", err)
	}

	t.Run("GetFileLink", func(t *testing.T) {
		link, err := c.GetFileLink(ctx, meta.FileID)
		if err != nil {
			t.Fatalf("get file link failed: %v", err)
		}
		if len(link.Hosts) == 0 {
			t.Fatal("no hosts in response")
		}
		if link.Path == "" {
			t.Fatal("path is empty")
		}
		url := link.URL()
		if url == "" {
			t.Fatal("url is empty")
		}
	})

	t.Run("GetFileLinkWithOpts", func(t *testing.T) {
		link, err := c.GetFileLinkWithOpts(ctx, meta.FileID, &pcloud.FileLinkOpts{
			ForceDownload: true,
		})
		if err != nil {
			t.Fatalf("get file link with opts failed: %v", err)
		}
		if link.URL() == "" {
			t.Fatal("url is empty")
		}
	})
}

func TestWalk(t *testing.T) {
	c, ctx := getClient(t)
	defer c.Logout(ctx)

	testFolder := "pcloud_walk_test_" + time.Now().Format("20060102150405")
	folder, err := c.CreateFolder(ctx, 0, testFolder)
	if err != nil {
		t.Fatalf("create test folder failed: %v", err)
	}
	defer c.DeleteFolderRecursive(ctx, folder.FolderID)

	c.CreateFolder(ctx, folder.FolderID, "sub1")
	sub2, _ := c.CreateFolder(ctx, folder.FolderID, "sub2")
	c.CreateFolder(ctx, sub2.FolderID, "sub2_nested")
	c.Upload(ctx, folder.FolderID, "file1.txt", bytes.NewReader([]byte("1")), nil)
	c.Upload(ctx, sub2.FolderID, "file2.txt", bytes.NewReader([]byte("2")), nil)

	t.Run("Walk", func(t *testing.T) {
		var items []string
		for item, err := range c.Walk(ctx, folder.FolderID) {
			if err != nil {
				t.Fatalf("walk failed: %v", err)
			}
			items = append(items, item.Name)
		}

		expected := []string{"sub1", "sub2", "sub2_nested", "file2.txt", "file1.txt"}
		if len(items) != len(expected) {
			t.Fatalf("expected %d items, got %d: %v", len(expected), len(items), items)
		}

		expectedSet := make(map[string]bool)
		for _, e := range expected {
			expectedSet[e] = true
		}
		for _, item := range items {
			if !expectedSet[item] {
				t.Fatalf("unexpected item: %s", item)
			}
		}
	})

	t.Run("WalkEarlyBreak", func(t *testing.T) {
		count := 0
		for _, err := range c.Walk(ctx, folder.FolderID) {
			if err != nil {
				t.Fatalf("walk failed: %v", err)
			}
			count++
			if count >= 2 {
				break
			}
		}
		if count != 2 {
			t.Fatalf("expected 2 items before break, got %d", count)
		}
	})
}
