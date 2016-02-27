# https://docs.pcloud.com


```
package main

import "github.com/yanmhlv/pcloud"

func main() {
    c := NewClient()
    fmt.Println("Login", c.Login("myemail", "mypassword"))
    fmt.Println("CreateFolder", c.CreateFolder("/helloworld", 0, ""))
    fmt.Println("DeleteFolder", c.DeleteFolder("/helloworld", 0))
    fmt.Println("DeleteFolderRecursive", c.DeleteFolderRecursive("/1/2/3/4", 0))
    fmt.Println(c.Auth, c.Client)
    fmt.Println(c.Logout())
}
```
