package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	data := "Hello, OTUS!"

	fmt.Println(stringutil.Reverse(data))
}
