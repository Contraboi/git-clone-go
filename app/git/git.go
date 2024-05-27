package git

type Git struct {
	name         string
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

func NewGit(name string) Git {
	master := Branch{
		Name:   "master",
		Commit: nil,
	}

	return Git{name: name, Head: master, lastCommitId: 0}
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
