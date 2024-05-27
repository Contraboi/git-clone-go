package git

import (
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"strings"
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

func Init(name string) Git {
	master := Branch{
		Name:   "master",
		Commit: nil,
	}

	for _, dir := range []string{".git-clone", ".git-clone/objects", ".git-clone/refs"} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			os.Exit(1)
		}
	}

	headFileContents := []byte("ref: refs/heads/master\n")
	if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Initialized git directory")
	return Git{Name: name, Head: master}
}

func CatFile(sha string) {
	path := fmt.Sprintf(".git/objects/%v/%v", sha[0:2], sha[2:])
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err)
		os.Exit(1)
	}

	r, err := zlib.NewReader(io.Reader(file))
	defer r.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating zlib reader: %s\n", err)
		os.Exit(1)
	}
	s, err := io.ReadAll(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading zlib reader: %s\n", err)
		os.Exit(1)
	}
	parts := strings.Split(string(s), "\x00")
	fmt.Print(parts[1])
}

func (g *Git) Commit(message string) Commit {
	cmt := Commit{
		Id:      sha1.New().Sum([]byte(message)),
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
