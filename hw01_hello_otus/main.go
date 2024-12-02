package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	helloString := "Hello, OTUS!"
	reverHelloString := reverse.String(helloString)

	fmt.Println(reverHelloString)
}
