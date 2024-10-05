package main

import (
	"fmt"
	"os"
)

const minArgs = 2

func main() {
	args := os.Args[1:]
	if len(args) < minArgs {
		fmt.Println("low args")
		return
	}

	dir, cmd := args[0], args[1:]
	envs, err := ReadDir(dir)
	if err != nil {
		fmt.Printf("err %v", err)
		return
	}

	exit := RunCmd(cmd, envs)
	os.Exit(exit)
}
