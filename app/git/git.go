package git

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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

func getPath(sha string) string {
	return fmt.Sprintf(".git-clone/objects/%v/%v", sha[0:2], sha[2:])
}

func CatFile(sha string) {
	file, err := os.Open(getPath(sha))
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

func HashObject(fileName string) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(1)
	}
	content := string(file)
	contentAndHeader := fmt.Sprintf("blob %d\x00%s", len(file), content)
	sha := (sha1.Sum([]byte(contentAndHeader)))
	hash := fmt.Sprintf("%x", sha)
	blobName := []rune(hash)
	blobPath := ".git-clone/objects/"

	// TODO: Think about this a bit
	for i, v := range blobName {
		blobPath += string(v)
		if i == 1 {
			blobPath += "/"
		}
	}

	var buffer bytes.Buffer

	z := zlib.NewWriter(&buffer)
	z.Write([]byte(contentAndHeader))
	z.Close()

	os.MkdirAll(filepath.Dir(blobPath), os.ModePerm)

	f, _ := os.Create(blobPath)
	defer f.Close()

	f.Write(buffer.Bytes())
	fmt.Print(hash)
}

func LsTree(sha string) {
	if len(sha) != 40 {
		log.Fatalf("Invalid SHA-1 hash length: %s\n", sha)
	}

	fmt.Println(sha[:2], sha[2:])
	file, err := os.Open(fmt.Sprintf(".git/objects/%s/%s", sha[:2], sha[2:]))
	if err != nil {
		log.Fatalf("Error opening file: %s\n", err)
	}
	defer file.Close()

	zlibReader, err := zlib.NewReader(file)
	if err != nil {
		log.Fatalf("Error creating zlib reader: %s\n", err)
	}
	defer zlibReader.Close()

	body, err := io.ReadAll(zlibReader)
	if err != nil {
		log.Fatalf("Error reading zlib reader: %s\n", err)
	}
	files := []string{}
	body = bytes.SplitN(body, []byte("\x00"), 2)[1]

	fmt.Println(string(body), " - ", sha)
	// TODO: Parse the tree object
	for _, line := range bytes.Split(body, []byte("\n")) {
		if len(line) == 0 {
			continue
		}

		parts := bytes.SplitN(line, []byte(" "), 2)
		if string(parts[0]) == "tree" {
			files = append(files, string(parts[1]))

			LsTree(string(parts[1]))
		}

	}

	fmt.Println(files)

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
