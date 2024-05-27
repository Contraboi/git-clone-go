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
		g, err := git.Init(name)
		if err != nil {
			os.Exit(1)
		}

		fmt.Println("Initialized git repository", g.Name)
	default:
		fmt.Println("Unknown command")
	}
}
