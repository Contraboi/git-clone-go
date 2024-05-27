package git

type Git struct {
	name         string
	lastCommitId int
}

func NewGit(name string) Git {
	return Git{name: name, lastCommitId: -1}
}

func (g Git) Commit(message string) Commit {
	g.lastCommitId++
	return newCommit(g.lastCommitId, message)
}
