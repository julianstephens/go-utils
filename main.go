package main

import (
	"fmt"

	"github.com/julianstephens/go-utils/httputil/auth"
)


func main() {
	pwd, _ := auth.HashPassword("test")
	fmt.Println(pwd)
}