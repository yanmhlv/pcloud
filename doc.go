// Package pcloud provides a Go client for the pCloud API.
//
// # Authentication
//
// Create a client and authenticate with username/password:
//
//	c := pcloud.NewClient(pcloud.BaseURLUS)  // or BaseURLEU for EU
//	err := c.Login("user@example.com", "password")
//	defer c.Logout()
//
// Alternatively, set an existing auth token directly:
//
//	c := pcloud.NewClient(pcloud.BaseURLUS)
//	c.SetAuth("existing-token")
//
// # Folders
//
// List, create, rename, copy, and delete folders:
//
//	folder, _ := c.ListFolder(0, nil)  // 0 is root folder
//	c.CreateFolder(0, "new-folder")
//	c.RenameFolder(folderID, "new-name")
//	c.DeleteFolderRecursive(folderID)
//
// # Files
//
// Upload, download, and manage files:
//
//	meta, _ := c.Upload(folderID, "file.txt", reader, nil)
//	body, _ := c.Download(fileID)
//	c.RenameFile(fileID, "new-name.txt")
//	c.DeleteFile(fileID)
//
// # Streaming
//
// Get direct download links for files:
//
//	link, _ := c.GetFileLink(fileID)
//	url := link.URL()
//
// # Revisions
//
// List and revert file revisions:
//
//	revisions, _ := c.ListRevisions(fileID)
//	c.RevertRevision(fileID, revisionID)
//
// # Walking
//
// Recursively iterate over all files and folders using iter.Seq2:
//
//	for item, err := range c.Walk(0) {
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    fmt.Println(item.Path)
//	}
package pcloud
