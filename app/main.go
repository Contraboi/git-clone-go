package main

import (
	"fmt"
	"git-clone/app/git"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a command")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		if len(os.Args) < 3 {
			fmt.Println("usage: init <name>")
			os.Exit(1)
		}

		name := os.Args[2]
		git.Init(name)
	case "cat-file":
		if len(os.Args) < 3 {
			fmt.Println("usage: cat-file <sha>")
			os.Exit(1)
		}

		sha := os.Args[2]
		git.CatFile(sha)
	default:
		fmt.Println("Unknown command")
	}
}
