package git

import (
	"fmt"
	"os"
)

type Git struct {
	Name         string
	Head         Branch
	lastCommitId int
}

type Commit struct {
	Id      int
	Message string
	Parent  *Commit
}

type Branch struct {
	Name   string
	Commit *Commit
}

func Init(name string) (Git, error) {
	master := Branch{
		Name:   "master",
		Commit: nil,
	}

	err := os.Mkdir(".git-clone", 0755)
	if err != nil {
		fmt.Println(err)
		return Git{}, err
	}

	err = os.Mkdir(".git-clone/objects", 0755)
	if err != nil {
		fmt.Println(err)
		return Git{}, err

	}
	err = os.Mkdir(".git-clone/refs", 0755)
	if err != nil {
		fmt.Println(err)
		return Git{}, err
	}
	_, err = os.Create(".git-clone/HEAD")
	if err != nil {
		fmt.Println(err)
		return Git{}, err
	}

	return Git{Name: name, Head: master, lastCommitId: 0}, nil
}

func (g *Git) Commit(message string) Commit {
	cmt := Commit{
		Id:      g.lastCommitId,
		Message: message,
		Parent:  g.Head.Commit,
	}

	g.lastCommitId = g.lastCommitId + 1
	g.Head.Commit = &cmt

	return cmt
}

func (g *Git) Log() []Commit {
	cmt := g.Head.Commit
	history := []Commit{}

	for cmt.Parent != nil {
		history = append(history, *cmt)
		cmt = cmt.Parent
	}

	return history
}
