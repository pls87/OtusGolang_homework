package main

import (
	"os"
)

func main() {
	env, _ := ReadDir(os.Args[1])
	code := RunCmd(os.Args[2:], env)
	os.Exit(code)
}
