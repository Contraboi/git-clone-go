package git

import (
	"crypto/sha1"
	"fmt"
	"os"
)

type Git struct {
	Name string
	Head Branch
}

type Commit struct {
	Id      []byte
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

	dirName := ".git-clone"

	err := os.Mkdir(dirName, 0755)
	if err != nil {
		fmt.Println(err)
		return Git{}, err
	}

	err = os.Mkdir(dirName+"/objects", 0755)
	if err != nil {
		fmt.Println(err)
		return Git{}, err

	}
	err = os.Mkdir(dirName+"/refs", 0755)
	if err != nil {
		fmt.Println(err)
		return Git{}, err
	}
	file, err := os.Create(dirName + "/HEAD")
	if err != nil {
		fmt.Println(err)
		return Git{}, err
	}
	defer file.Close()

	_, err = file.WriteString("ref: refs/heads/master\n")
	if err != nil {
		fmt.Println(err)
		return Git{}, err
	}

	return Git{Name: name, Head: master}, nil
}

func (g *Git) Commit(message string) Commit {
	sha1 := sha1.New()
	id := sha1.Sum([]byte(message))

	cmt := Commit{
		Id:      id,
		Message: message,
		Parent:  g.Head.Commit,
	}

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
