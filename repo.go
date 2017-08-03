package main

import (
	"strings"
	"gopkg.in/src-d/go-git.v4"
)

type Repo struct {
	Path         string
	RemoteStatus string
	LocalStatus  string
	Repository   *git.Repository
}

func NewRepo(path string) *Repo {
	r, _ := git.PlainOpen(path)
	repo := &Repo{
		Path: path,
		Repository: r,
	}
	return repo
}

func (repo *Repo) Name() string {
	return repo.Path[strings.LastIndex(repo.Path, "/")+1:]
}
