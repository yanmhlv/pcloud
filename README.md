# https://docs.pcloud.com

```bash
username=myusername password=mypassword go test github.com/yanmhlv/pcloud
```

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
    fmt.Println("CreateFolder /helloworld", c.CreateFolder("/helloworld", 0, ""))
    fmt.Println("CreateFolder /helloworld/1", c.CreateFolder("/helloworld/1", 0, ""))
    fmt.Println("CreateFolder /helloworld/2", c.CreateFolder("/helloworld/2", 0, ""))

    fmt.Println("DeleteFolder /helloworld/2", c.DeleteFolder("/helloworld/2", 0))

    fmt.Println("RenameFolder /helloworld to /hello_world", c.RenameFolder(-1, "/helloworld", "/hello_world"))
    fmt.Println("DeleteFolderRecursive /hello_world", c.DeleteFolderRecursive("/hello_world", 0))

    fh, _ := os.Open("/Users/yan/Desktop/index.html")
    fmt.Println("UploadFile index.html", c.UploadFile(fh, "", 0, "index.html", 0, "", 0))
    fmt.Println("CopyFile /index.html to /index2.html", c.CopyFile(0, "/index.html", 0, "", "/index2.html"))
    fmt.Println("DeleteFile /index2.html", c.DeleteFile(0, "/index2.html"))
    fmt.Println("RenameFile /index.html to /index2.html", c.RenameFile(0, "/index.html", "/index2.html", 0, ""))
    fmt.Println("DeleteFile /index2.html", c.DeleteFile(0, "/index2.html"))

    fmt.Println("authkey:", *c.Auth)
    fmt.Println("Logout", c.Logout())
}

```
