package tests

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/yanmhlv/pcloud"
)

func getClient(t *testing.T) *pcloud.Client {
	t.Helper()

	username := os.Getenv("PCLOUD_USERNAME")
	password := os.Getenv("PCLOUD_PASSWORD")
	if username == "" || password == "" {
		t.Skip("PCLOUD_USERNAME and PCLOUD_PASSWORD env vars required")
	}

	baseURL := os.Getenv("PCLOUD_BASE_URL")
	c := pcloud.NewClient(baseURL)

	if err := c.Login(username, password); err != nil {
		t.Fatalf("login failed: %v", err)
	}

	return c
}

func TestAuth(t *testing.T) {
	username := os.Getenv("PCLOUD_USERNAME")
	password := os.Getenv("PCLOUD_PASSWORD")
	if username == "" || password == "" {
		t.Skip("PCLOUD_USERNAME and PCLOUD_PASSWORD env vars required")
	}

	baseURL := os.Getenv("PCLOUD_BASE_URL")
	c := pcloud.NewClient(baseURL)

	t.Run("Login", func(t *testing.T) {
		if err := c.Login(username, password); err != nil {
			t.Fatalf("login failed: %v", err)
		}
		if c.Auth() == "" {
			t.Fatal("auth token is empty after login")
		}
	})

	t.Run("UserInfo", func(t *testing.T) {
		info, err := c.UserInfo()
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
		if err := c.Logout(); err != nil {
			t.Fatalf("logout failed: %v", err)
		}
		if c.Auth() != "" {
			t.Fatal("auth token should be empty after logout")
		}
	})
}

func TestFolders(t *testing.T) {
	c := getClient(t)
	defer c.Logout()

	testFolder := "pcloud_test_" + time.Now().Format("20060102150405")
	var folderID uint64

	t.Run("CreateFolder", func(t *testing.T) {
		folder, err := c.CreateFolder(0, testFolder)
		if err != nil {
			t.Fatalf("create folder failed: %v", err)
		}
		if folder.Name != testFolder {
			t.Fatalf("expected name %s, got %s", testFolder, folder.Name)
		}
		folderID = folder.FolderID
	})

	t.Run("ListFolder", func(t *testing.T) {
		folder, err := c.ListFolder(0, nil)
		if err != nil {
			t.Fatalf("list folder failed: %v", err)
		}
		if !folder.IsFolder {
			t.Fatal("root should be a folder")
		}

		found := false
		for _, item := range folder.Contents {
			if item.FolderID == folderID {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("created folder not found in root")
		}
	})

	t.Run("CreateFolderIfNotExists", func(t *testing.T) {
		folder, err := c.CreateFolderIfNotExists(folderID, "subfolder")
		if err != nil {
			t.Fatalf("create subfolder failed: %v", err)
		}
		if folder.Name != "subfolder" {
			t.Fatalf("expected name subfolder, got %s", folder.Name)
		}

		folder2, err := c.CreateFolderIfNotExists(folderID, "subfolder")
		if err != nil {
			t.Fatalf("create existing subfolder failed: %v", err)
		}
		if folder2.FolderID != folder.FolderID {
			t.Fatal("should return same folder")
		}
	})

	t.Run("RenameFolder", func(t *testing.T) {
		newName := testFolder + "_renamed"
		folder, err := c.RenameFolder(folderID, newName)
		if err != nil {
			t.Fatalf("rename folder failed: %v", err)
		}
		if folder.Name != newName {
			t.Fatalf("expected name %s, got %s", newName, folder.Name)
		}

		folder, _ = c.RenameFolder(folderID, testFolder)
		if folder.Name != testFolder {
			t.Fatal("rename back failed")
		}
	})

	t.Run("CopyFolder", func(t *testing.T) {
		copyName := testFolder + "_copy"
		copyFolder, err := c.CreateFolder(0, copyName)
		if err != nil {
			t.Fatalf("create copy target failed: %v", err)
		}
		defer c.DeleteFolderRecursive(copyFolder.FolderID)

		subFolder, _ := c.CreateFolder(folderID, "to_copy")
		copied, err := c.CopyFolder(subFolder.FolderID, copyFolder.FolderID)
		if err != nil {
			t.Fatalf("copy folder failed: %v", err)
		}
		if copied.Name != "to_copy" {
			t.Fatal("copied folder has wrong name")
		}
	})

	t.Run("DeleteFolderRecursive", func(t *testing.T) {
		if err := c.DeleteFolderRecursive(folderID); err != nil {
			t.Fatalf("delete folder recursive failed: %v", err)
		}

		folder, _ := c.ListFolder(0, nil)
		for _, item := range folder.Contents {
			if item.FolderID == folderID {
				t.Fatal("folder should be deleted")
			}
		}
	})
}

func TestFiles(t *testing.T) {
	c := getClient(t)
	defer c.Logout()

	testFolder := "pcloud_file_test_" + time.Now().Format("20060102150405")
	folder, err := c.CreateFolder(0, testFolder)
	if err != nil {
		t.Fatalf("create test folder failed: %v", err)
	}
	defer c.DeleteFolderRecursive(folder.FolderID)

	testContent := []byte("hello pcloud test content")
	var fileID uint64

	t.Run("Upload", func(t *testing.T) {
		meta, err := c.Upload(folder.FolderID, "test.txt", bytes.NewReader(testContent), nil)
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
		meta, err := c.Stat(fileID)
		if err != nil {
			t.Fatalf("stat failed: %v", err)
		}
		if meta.Name != "test.txt" {
			t.Fatalf("expected name test.txt, got %s", meta.Name)
		}
	})

	t.Run("Download", func(t *testing.T) {
		body, err := c.Download(fileID)
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
		meta, err := c.RenameFile(fileID, "renamed.txt")
		if err != nil {
			t.Fatalf("rename file failed: %v", err)
		}
		if meta.Name != "renamed.txt" {
			t.Fatalf("expected name renamed.txt, got %s", meta.Name)
		}
	})

	t.Run("CopyFile", func(t *testing.T) {
		meta, err := c.CopyFile(fileID, folder.FolderID)
		if err != nil {
			t.Fatalf("copy file failed: %v", err)
		}
		if meta.FileID == fileID {
			t.Fatal("copied file should have different id")
		}
		c.DeleteFile(meta.FileID)
	})

	t.Run("DeleteFile", func(t *testing.T) {
		if err := c.DeleteFile(fileID); err != nil {
			t.Fatalf("delete file failed: %v", err)
		}

		_, err := c.Stat(fileID)
		if err == nil {
			t.Fatal("stat should fail for deleted file")
		}
	})
}

func TestRevisions(t *testing.T) {
	c := getClient(t)
	defer c.Logout()

	testFolder := "pcloud_rev_test_" + time.Now().Format("20060102150405")
	folder, err := c.CreateFolder(0, testFolder)
	if err != nil {
		t.Fatalf("create test folder failed: %v", err)
	}
	defer c.DeleteFolderRecursive(folder.FolderID)

	content1 := []byte("version 1")
	content2 := []byte("version 2")

	meta, err := c.Upload(folder.FolderID, "revtest.txt", bytes.NewReader(content1), nil)
	if err != nil {
		t.Fatalf("upload v1 failed: %v", err)
	}
	fileID := meta.FileID

	_, err = c.Upload(folder.FolderID, "revtest.txt", bytes.NewReader(content2), nil)
	if err != nil {
		t.Fatalf("upload v2 failed: %v", err)
	}

	t.Run("ListRevisions", func(t *testing.T) {
		revs, err := c.ListRevisions(fileID)
		if err != nil {
			t.Fatalf("list revisions failed: %v", err)
		}
		if len(revs) < 2 {
			t.Fatalf("expected at least 2 revisions, got %d", len(revs))
		}
	})

	t.Run("RevertRevision", func(t *testing.T) {
		revs, _ := c.ListRevisions(fileID)
		if len(revs) < 2 {
			t.Skip("not enough revisions")
		}

		oldRev := revs[len(revs)-1]
		_, err := c.RevertRevision(fileID, oldRev.RevisionID)
		if err != nil {
			t.Fatalf("revert revision failed: %v", err)
		}

		body, _ := c.Download(fileID)
		defer body.Close()
		content, _ := io.ReadAll(body)
		if !bytes.Equal(content, content1) {
			t.Fatalf("content should be reverted to v1, got %s", content)
		}
	})
}

func TestStreaming(t *testing.T) {
	c := getClient(t)
	defer c.Logout()

	testFolder := "pcloud_stream_test_" + time.Now().Format("20060102150405")
	folder, err := c.CreateFolder(0, testFolder)
	if err != nil {
		t.Fatalf("create test folder failed: %v", err)
	}
	defer c.DeleteFolderRecursive(folder.FolderID)

	content := []byte("stream test content")
	meta, err := c.Upload(folder.FolderID, "stream.txt", bytes.NewReader(content), nil)
	if err != nil {
		t.Fatalf("upload failed: %v", err)
	}

	t.Run("GetFileLink", func(t *testing.T) {
		link, err := c.GetFileLink(meta.FileID)
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
		link, err := c.GetFileLinkWithOpts(meta.FileID, &pcloud.FileLinkOpts{
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
	c := getClient(t)
	defer c.Logout()

	testFolder := "pcloud_walk_test_" + time.Now().Format("20060102150405")
	folder, err := c.CreateFolder(0, testFolder)
	if err != nil {
		t.Fatalf("create test folder failed: %v", err)
	}
	defer c.DeleteFolderRecursive(folder.FolderID)

	c.CreateFolder(folder.FolderID, "sub1")
	sub2, _ := c.CreateFolder(folder.FolderID, "sub2")
	c.CreateFolder(sub2.FolderID, "sub2_nested")
	c.Upload(folder.FolderID, "file1.txt", bytes.NewReader([]byte("1")), nil)
	c.Upload(sub2.FolderID, "file2.txt", bytes.NewReader([]byte("2")), nil)

	t.Run("Walk", func(t *testing.T) {
		var items []string
		for item, err := range c.Walk(folder.FolderID) {
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
		for _, err := range c.Walk(folder.FolderID) {
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
