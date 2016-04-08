package pcloud

import (
	"bytes"
	"os"
	"testing"
)

var client *pCloudClient

func init() {
	client = NewClient()
}

// TestLogin
func TestLogin(t *testing.T) {
	if err := client.Login(os.Getenv("username"), os.Getenv("password")); err != nil {
		t.Fatal("can't login", err)
	}
}

// TestAuthkey
func TestAuthkey(t *testing.T) {
	if client.Auth == nil {
		t.Error("auth key is nil!")
	}
}

// TestCreateFolder
func TestCreateFolder(t *testing.T) {
	if err := client.CreateFolder("/helloworld", -1, ""); err != nil {
		t.Error("create folder error; path: /helloworld", err)
	}
	if err := client.CreateFolder("/helloworld/1", -1, ""); err != nil {
		t.Error("create embedded folder error; path: /helloworld/1", err)
	}
	if err := client.CreateFolder("/helloworld/2", -1, ""); err != nil {
		t.Error("create embedded folder error; path: /helloworld/2", err)
	}
	if err := client.CreateFolder("", 0, "testfolder"); err != nil {
		t.Error("create embedded folder by id; /testfolder", err)
	}
}

// TestDeleteFolder
func TestDeleteFolder(t *testing.T) {
	if err := client.DeleteFolder("/helloworld/2", -1); err != nil {
		t.Error("delete folder error; path: /helloworld/2", err)
	}
}

// TestDeleteFolderRecursive
func TestDeleteFolderRecursive(t *testing.T) {
	if err := client.DeleteFolderRecursive("/hello_world", -1); err != nil {
		t.Error("delete folder recursive error; path: /hello_world", err)
	}
}

// TestRenameFolder
func TestRenameFolder(t *testing.T) {
	if err := client.RenameFolder(-1, "/helloworld", "/hello_world"); err != nil {
		t.Error("rename folder error; rename /helloworld to /hello_world", err)
	}
}

// TestUploadFile
func TestUploadFile(t *testing.T) {
	buf := bytes.NewBuffer([]byte("test data"))
	if err := client.UploadFile(buf, "", 0, "testfile", 0, "", 0); err != nil {
		t.Error("upload testfile error", err)
	}
}

// TestCopyFile
func TestCopyFile(t *testing.T) {
	if err := client.CopyFile(0, "/testfile", 0, "", "/testfile_copy"); err != nil {
		t.Error("copy testfile error", err)
	}
}

// TestDeleteFile
func TestDeleteFile(t *testing.T) {
	if err := client.DeleteFile(0, "/testfile_copy"); err != nil {
		t.Error("delete file error", err)
	}
}

// TestRenameFile
func TestRenameFile(t *testing.T) {
	if err := client.RenameFile(0, "/testfile", "/testfile_rename", 0, ""); err != nil {
		t.Error("rename file error", err)
	}
}

// TestLogout
func TestLogout(t *testing.T) {
	if err := client.Logout(); err != nil {
		t.Error("logout error", err)
	}
}
