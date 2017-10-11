package main

import (
	"strings"
)

type Repo struct {
	Path         string
	RepoUrl      string
	RemoteStatus string
	LocalStatus  string
}

func (repo *Repo) Name() string {
	return repo.Path[strings.LastIndex(repo.Path, "/")+1:]
}
