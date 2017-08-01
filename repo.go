package main

import (
	"strings"
	"fmt"
)

type Repo struct {
	Path string
	RemoteStatus string
	LocalStatus string
}

func (repo *Repo) Name() string {
	return repo.Path[strings.LastIndex(repo.Path, "/")+1:]
}

func (repo *Repo) Fetch() {
	run(repo.Path, []string{"fetch", "--all"})
}

func (repo *Repo) Status() string {
	output, err := run(repo.Path, []string{"status", "-unormal"})
	if err != nil || strings.Contains(output, "Not a git repo") {
		return "error"
	}
	return remoteStatus(output) + " " + localStatus(output)
}

func (repo *Repo) Branch() string {
	output, err := run(repo.Path, []string{"rev-parse", "--abbrev-ref", "HEAD"})
	if err != nil {
		fmt.Println(err)
	}
	return strings.TrimSuffix(output, "\n")
}

func (repo *Repo) Branches(currentBranch string) string {
	output, err := run(repo.Path, []string{"branch", "--column", "--format=%(refname:short)~"})
	if err != nil {
		fmt.Println(err)
	}
	output = strings.Replace(output, "\n", "", -1)
	output = strings.Replace(output, currentBranch+"~", "", -1)
	output = strings.Replace(output, " ", "", -1)
	output = strings.Replace(output, "~", " ", -1)
	return output
}


func localStatus(output string) string {
	if strings.Contains(output, "nothing added to commit but untracked files present") {
		return UNTRACKED
	} else if strings.Contains(output, "nothing to commit") {
		return SYNC
	}

	return CHANGED
}
func remoteStatus(output string) string {
	if strings.Contains(output, "have diverged") {
		return DIVERGED
	} else if strings.Contains(output, "branch is ahead") {
		return AHEAD
	} else if strings.Contains(output, "branch is behind") {
		return BEHIND
	} else if strings.Contains(output, "branch is up-to-date") {
		return SYNC
	}
	return NO_REMOTE
}


