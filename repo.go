package main

import (
	"strings"
)

type Repo struct {
	Path         string
	RemoteStatus string
	LocalStatus  string
}

func (repo *Repo) Name() string {
	return repo.Path[strings.LastIndex(repo.Path, "/")+1:]
}
