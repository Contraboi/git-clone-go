package main

import (
	"fmt"
	"git-clone/app/git"
)

func main() {
	g := git.NewGit("my-repo")
	g.Commit("commit 1")
	g.Commit("commit 2")
	g.Commit("commit 3")
	g.Commit("commit 4")
	g.Commit("commit 6")

	for _, cmt := range g.Log() {
		fmt.Println("Commit ID:", cmt.Id, "Message:", cmt.Message)
	}
}
