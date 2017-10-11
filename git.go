package main

import (
	"fmt"
	"strings"
)

type Git interface {
	Clone(repo *Repo) string

	Fetch(repo *Repo)

	Status(repo *Repo) string

	Branch(repo *Repo) string

	Branches(repo *Repo, currentBranch string) string
}

type CommandLineGit struct {
	runner CommandLineRunner
}

func NewDefaultGit() Git {
	return NewCustomGit(&DefaultCommandLineRunner{})
}

func NewCustomGit(runner CommandLineRunner) Git {
	git := &CommandLineGit{
		runner: runner,
	}
	return git
}

func (git *CommandLineGit) Clone(repo *Repo) string {
	output, err := git.runner.Run("", "git", []string{"clone", repo.RepoUrl, repo.Path})
	if err != nil || strings.Contains(output, "Not a git repo") {
		return "error"
	}
	return "ok"
}

func (git *CommandLineGit) Fetch(repo *Repo) {
	git.runner.Run(repo.Path, "git", []string{"fetch", "--all"})
}

func (git *CommandLineGit) Status(repo *Repo) string {
	output, err := git.runner.Run(repo.Path, "git", []string{"status", "-unormal"})
	if err != nil || strings.Contains(output, "Not a git repo") {
		return "error"
	}
	return remoteStatus(output) + " " + localStatus(output)
}

func (git *CommandLineGit) Branch(repo *Repo) string {
	output, err := git.runner.Run(repo.Path, "git", []string{"rev-parse", "--abbrev-ref", "HEAD"})
	if err != nil {
		fmt.Println(err)
	}
	return strings.TrimSuffix(output, "\n")
}

func (git *CommandLineGit) Branches(repo *Repo, currentBranch string) string {
	output, err := git.runner.Run(repo.Path, "git", []string{"branch", "--column", "--format=%(refname:short)~"})
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
