# https://docs.pcloud.com


```go
package main

import (
    "fmt"
    "os"

    "github.com/yanmhlv/pcloud"
)

func main() {
    c := pcloud.NewClient()
    fmt.Println("Login", c.Login("myemail", "mypassword"))
    fmt.Println("CreateFolder", c.CreateFolder("/helloworld", 0, ""))
    fmt.Println("DeleteFolder", c.DeleteFolder("/helloworld", 0))
    fmt.Println("DeleteFolderRecursive", c.DeleteFolderRecursive("/1/2/3/4", 0))

    fmt.Println("create folder 2", c.CreateFolder("/test", 0, ""))
    fmt.Println("RenameFolder", c.RenameFolder(0, "/test", "/hello_world"))
    c.DeleteFolder("/hello_world", 0)
    fmt.Println(c.Auth, c.Client)
    fmt.Println("CopyFile", c.CopyFile(0, "/1.go", 0, "", "/2.go"))
    fmt.Println("RenameFile", c.RenameFile(0, "/2.go", "/3.go", 0, ""))
    fmt.Println("DeleteFile", c.DeleteFile(0, "/3.go"))

    fh, _ := os.Open("/Users/yan/Desktop/index.html")
    fmt.Println("UploadFile", c.UploadFile(fh, "", 0, "test.go", 0, "", 0))
    fmt.Println("Logout", c.Logout())
}

```
