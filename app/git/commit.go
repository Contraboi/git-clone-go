package git

type Commit struct {
	// TODO: SHA myb
	Id      int
	Message string
}

func newCommit(id int, message string) Commit {
	return Commit{Id: id, Message: message}
}
