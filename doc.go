// Package pcloud provides a Go client for the pCloud API.
//
// # Authentication
//
// Create a client and authenticate with username/password:
//
//	ctx := context.Background()
//	c := pcloud.NewClient(pcloud.BaseURLUS)  // or BaseURLEU for EU
//	err := c.Login(ctx, "user@example.com", "password")
//	defer c.Logout(ctx)
//
// Alternatively, use an OAuth2 token:
//
//	ctx := context.Background()
//	c := pcloud.NewClient(pcloud.BaseURLUS)
//	c.SetTokenSource(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "token"}))
//
// # Folders
//
// List, create, rename, copy, and delete folders:
//
//	ctx := context.Background()
//	folder, _ := c.ListFolder(ctx, 0, nil)  // 0 is root folder
//	c.CreateFolder(ctx, 0, "new-folder")
//	c.RenameFolder(ctx, folderID, "new-name")
//	c.DeleteFolderRecursive(ctx, folderID)
//
// # Files
//
// Upload, download, and manage files:
//
//	ctx := context.Background()
//	meta, _ := c.Upload(ctx, folderID, "file.txt", reader, nil)
//	body, _ := c.Download(ctx, fileID, nil)
//	c.RenameFile(ctx, fileID, "new-name.txt")
//	c.DeleteFile(ctx, fileID)
//
// # Streaming
//
// Get direct download links for files:
//
//	ctx := context.Background()
//	link, _ := c.GetFileLink(ctx, fileID)
//	url := link.URL()
//
// # Revisions
//
// List and revert file revisions:
//
//	ctx := context.Background()
//	revisions, _ := c.ListRevisions(ctx, fileID)
//	c.RevertRevision(ctx, fileID, revisionID)
//
// # Walking
//
// Recursively iterate over all files and folders using iter.Seq2:
//
//	ctx := context.Background()
//	for item, err := range c.Walk(ctx, 0) {
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    fmt.Println(item.Path)
//	}
package pcloud
