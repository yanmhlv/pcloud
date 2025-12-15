package pcloud_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/yanmhlv/pcloud"
)

func Example() {
	ctx := context.Background()
	c := pcloud.NewClient(pcloud.BaseURLUS)

	if err := c.Login(ctx, "user@example.com", "password"); err != nil {
		log.Fatal(err)
	}
	defer c.Logout(ctx)

	folder, err := c.ListFolder(ctx, 0, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range folder.Contents {
		fmt.Println(item.Name)
	}
}

func ExampleClient_Login() {
	ctx := context.Background()
	c := pcloud.NewClient(pcloud.BaseURLUS)

	if err := c.Login(ctx, "user@example.com", "password"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("logged in successfully")
}

func ExampleClient_UserInfo() {
	ctx := context.Background()
	c := pcloud.NewClient(pcloud.BaseURLUS)
	c.Login(ctx, "user@example.com", "password")
	defer c.Logout(ctx)

	info, err := c.UserInfo(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Email: %s\n", info.Email)
	fmt.Printf("Quota: %d bytes\n", info.Quota)
	fmt.Printf("Used: %d bytes\n", info.UsedQuota)
}

func ExampleClient_ListFolder() {
	ctx := context.Background()
	c := pcloud.NewClient(pcloud.BaseURLUS)
	c.Login(ctx, "user@example.com", "password")
	defer c.Logout(ctx)

	folder, err := c.ListFolder(ctx, 0, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range folder.Contents {
		if item.IsFolder {
			fmt.Printf("[DIR]  %s\n", item.Name)
		} else {
			fmt.Printf("[FILE] %s (%d bytes)\n", item.Name, item.Size)
		}
	}
}

func ExampleClient_ListFolder_recursive() {
	ctx := context.Background()
	c := pcloud.NewClient(pcloud.BaseURLUS)
	c.Login(ctx, "user@example.com", "password")
	defer c.Logout(ctx)

	folder, err := c.ListFolder(ctx, 0, &pcloud.ListFolderOpts{
		Recursive: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	var walk func(items []pcloud.Metadata, indent string)
	walk = func(items []pcloud.Metadata, indent string) {
		for _, item := range items {
			fmt.Printf("%s%s\n", indent, item.Name)
			if item.IsFolder {
				walk(item.Contents, indent+"  ")
			}
		}
	}
	walk(folder.Contents, "")
}

func ExampleClient_Walk() {
	ctx := context.Background()
	c := pcloud.NewClient(pcloud.BaseURLUS)
	c.Login(ctx, "user@example.com", "password")
	defer c.Logout(ctx)

	for item, err := range c.Walk(ctx, 0) {
		if err != nil {
			log.Fatal(err)
		}
		if item.IsFolder {
			fmt.Printf("[DIR]  %s\n", item.Path)
		} else {
			fmt.Printf("[FILE] %s (%d bytes)\n", item.Path, item.Size)
		}
	}
}

func ExampleClient_Walk_findLargeFiles() {
	ctx := context.Background()
	c := pcloud.NewClient(pcloud.BaseURLUS)
	c.Login(ctx, "user@example.com", "password")
	defer c.Logout(ctx)

	const maxSize = 100 * 1024 * 1024 // 100MB
	for item, err := range c.Walk(ctx, 0) {
		if err != nil {
			log.Fatal(err)
		}
		if !item.IsFolder && item.Size > maxSize {
			fmt.Printf("%s: %d MB\n", item.Path, item.Size/1024/1024)
		}
	}
}

func ExampleClient_CreateFolder() {
	ctx := context.Background()
	c := pcloud.NewClient(pcloud.BaseURLUS)
	c.Login(ctx, "user@example.com", "password")
	defer c.Logout(ctx)

	folder, err := c.CreateFolder(ctx, 0, "my-new-folder")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created folder: %s (id: %d)\n", folder.Name, folder.FolderID)
}

func ExampleClient_Upload() {
	ctx := context.Background()
	c := pcloud.NewClient(pcloud.BaseURLUS)
	c.Login(ctx, "user@example.com", "password")
	defer c.Logout(ctx)

	content := bytes.NewReader([]byte("hello world"))
	meta, err := c.Upload(ctx, 0, "hello.txt", content, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Uploaded: %s (%d bytes)\n", meta.Name, meta.Size)
}

func ExampleClient_Upload_fromFile() {
	ctx := context.Background()
	c := pcloud.NewClient(pcloud.BaseURLUS)
	c.Login(ctx, "user@example.com", "password")
	defer c.Logout(ctx)

	f, err := os.Open("/path/to/local/file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	meta, err := c.Upload(ctx, 0, "file.txt", f, &pcloud.UploadOpts{
		RenameIfExists: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Uploaded: %s\n", meta.Name)
}

func ExampleClient_Download() {
	ctx := context.Background()
	c := pcloud.NewClient(pcloud.BaseURLUS)
	c.Login(ctx, "user@example.com", "password")
	defer c.Logout(ctx)

	body, err := c.Download(ctx, 12345)
	if err != nil {
		log.Fatal(err)
	}
	defer body.Close()

	content, err := io.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Downloaded %d bytes\n", len(content))
}

func ExampleClient_GetFileLink() {
	ctx := context.Background()
	c := pcloud.NewClient(pcloud.BaseURLUS)
	c.Login(ctx, "user@example.com", "password")
	defer c.Logout(ctx)

	link, err := c.GetFileLink(ctx, 12345)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Download URL: %s\n", link.URL())
	fmt.Printf("Expires: %s\n", link.Expires)
}

func ExampleClient_ListRevisions() {
	ctx := context.Background()
	c := pcloud.NewClient(pcloud.BaseURLUS)
	c.Login(ctx, "user@example.com", "password")
	defer c.Logout(ctx)

	revisions, err := c.ListRevisions(ctx, 12345)
	if err != nil {
		log.Fatal(err)
	}

	for _, rev := range revisions {
		fmt.Printf("Revision %d: %d bytes, created %s\n",
			rev.RevisionID, rev.Size, rev.Created)
	}
}

func ExampleClient_RevertRevision() {
	ctx := context.Background()
	c := pcloud.NewClient(pcloud.BaseURLUS)
	c.Login(ctx, "user@example.com", "password")
	defer c.Logout(ctx)

	revisions, err := c.ListRevisions(ctx, 12345)
	if err != nil {
		log.Fatal(err)
	}

	if len(revisions) < 2 {
		log.Fatal("no previous revisions")
	}

	oldRev := revisions[len(revisions)-1]
	meta, err := c.RevertRevision(ctx, 12345, oldRev.RevisionID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Reverted to revision, new size: %d\n", meta.Size)
}
