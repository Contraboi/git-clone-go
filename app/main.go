package main

import (
	"fmt"
	"git-clone/app/git"
)

func main() {
	g := git.NewGit("my-repo")
	c1 := g.Commit("first commit")

	fmt.Printf("Commit id: %d\n", c1.Id)
}
