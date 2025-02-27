package main

import (
	"fmt"
	"os"
)

func main() {
	osArgs := os.Args
	if len(osArgs) < 3 {
		fmt.Println("Wrong arguments")
		os.Exit(1)
	}

	dir := osArgs[1]
	args := osArgs[2:]

	env, err := ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}

	code := RunCmd(args, env)
	os.Exit(code)
}
