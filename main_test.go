package pcloud

import (
	"bytes"
	"math/rand"
	"os"
	"testing"
	"time"
)

func randomRune() rune {
	var asciiLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	return asciiLetters[rand.Intn(len(asciiLetters))]
}

func randomString(size int) string {
	var result string
	for i := 0; i < size; i++ {
		result += string(randomRune())
	}
	return result
}

var client *pCloudClient

var (
	folderByPath         string
	embeddedFolderByPath string
	folderNameByID       string
	folderByPathRename   string

	rootFodlerID         = 0
	beforeFilename       string
	beforeFilenameCopy   string
	beforeFilenameRename string
)

func init() {
	rand.Seed(time.Now().Unix())
	client = NewClient()

	folderByPath = "/" + randomString(20)
	embeddedFolderByPath = folderByPath + "/" + randomString(20)
	folderNameByID = randomString(20)
	folderByPathRename = folderByPath + "_rename"

	beforeFilename = randomString(20)
	beforeFilenameRename = beforeFilename + "_rename"
	beforeFilenameCopy = beforeFilename + "_copy"
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
func TestCreateFolderByPath(t *testing.T) {
	if err := client.CreateFolder(folderByPath, -1, ""); err != nil {
		t.Error("create folder error;", folderByPath, err)
	}
	if err := client.CreateFolder(folderByPath, -1, ""); err == nil {
		t.Error("duplicate create folder error;", folderByPath, err)
	}
}

// TestCreateEmbeddedFolderByPath
func TestCreateEmbeddedFolderByPath(t *testing.T) {
	if err := client.CreateFolder(embeddedFolderByPath, -1, ""); err != nil {
		t.Error("create embedded folder error;", embeddedFolderByPath, err)
	}
	if err := client.CreateFolder(embeddedFolderByPath, -1, ""); err == nil {
		t.Error("duplicate create embedded folder error;", embeddedFolderByPath)
	}
}

// TestCreatefolderNameByID
func TestCreatefolderNameByID(t *testing.T) {
	if err := client.CreateFolder("", rootFodlerID, folderNameByID); err != nil {
		t.Error("create embedded folder by id;", folderNameByID, err)
	}
	if err := client.CreateFolder("", rootFodlerID, folderNameByID); err == nil {
		t.Error("duplicate create embedded folder by id;", folderNameByID)
	}
}

// TestDeleteFolderByPath
func TestDeleteFolderByPath(t *testing.T) {
	if err := client.DeleteFolder(embeddedFolderByPath, -1); err != nil {
		t.Error("delete folder by path error;", embeddedFolderByPath, err)
	}
}

// TestRenameFolderByPath
func TestRenameFolderByPath(t *testing.T) {
	if err := client.RenameFolder(-1, folderByPath, folderByPathRename); err != nil {
		t.Error("rename folder error; rename /helloworld to /hello_world", err)
	}
	folderByPath = folderByPathRename
}

// TestDeleteFolderRecursiveByPath
func TestDeleteFolderRecursiveByPath(t *testing.T) {
	if err := client.DeleteFolderRecursive(folderByPath, -1); err != nil {
		t.Error("delete folder by path recursive error", folderByPath, err)
	}
}

// TestUploadFile
func TestUploadFile(t *testing.T) {
	buf := bytes.NewBuffer([]byte("test data"))
	if err := client.UploadFile(buf, "", rootFodlerID, beforeFilename, 0, "", 0); err != nil {
		t.Error("upload testfile error", beforeFilename, err)
	}
}

// TestCopyFile
func TestCopyFile(t *testing.T) {
	if err := client.CopyFile(0, "/"+beforeFilename, 0, "", "/"+beforeFilenameCopy); err != nil {
		t.Error("copy testfile error", err)
	}
}

// TestRenameFile
func TestRenameFile(t *testing.T) {
	if err := client.RenameFile(0, "/"+beforeFilename, "/"+beforeFilenameRename, 0, ""); err != nil {
		t.Error("rename file error", beforeFilename, beforeFilenameRename, err)
	}
	// beforeFilename = beforeFilenameRename
}

// TestDeleteFile
func TestDeleteFile(t *testing.T) {
	if err := client.DeleteFile(0, "/"+beforeFilenameCopy); err != nil {
		t.Error("delete file error", err)
	}
}

// TestLogout
func TestLogout(t *testing.T) {
	if err := client.Logout(); err != nil {
		t.Error("logout error", err)
	}
}
