package main

import (
	"fmt"
	"pcloud"
)

func main() {
	fmt.Println("auth example")
	c := pcloud.NewClient("myemail", "mypassword")
	fmt.Println("\t", c.Login())
	fmt.Println("\tAuthkey", c.Authkey)
	fmt.Println("\tUsername", c.Username)
	fmt.Println("\tPassword", c.Password)
	fmt.Println("######")
}
