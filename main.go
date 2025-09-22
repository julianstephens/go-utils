package main

import (
	"fmt"

	"github.com/julianstephens/go-utils/httputil/auth"
)

func main() {
	pwd := "test"
	res, err := auth.HashPassword(pwd)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
