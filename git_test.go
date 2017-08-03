package main

import (
	"testing"
	"strings"
)

func TestFetch(t *testing.T) {
	mockRunner := NewMockRunner("")
	git := NewCustomGit(mockRunner)

	repo := &Repo{Path: "/test/"}
	git.Fetch(repo)

	if mockRunner.folder != "/test/" {
		t.Errorf("Folder should be /test/")
	}

	if mockRunner.command!= "git" {
		t.Errorf("Command should be git")
	}

	if len(mockRunner.args) != 2 {
		t.Errorf("Args size should be 2")
	}

	if mockRunner.args[0] != "fetch" {
		t.Errorf("Args 0 should be fetch")
	}

	if mockRunner.args[1] != "--all" {
		t.Errorf("Args 1 should be --all")
	}
}

func TestStatus(t *testing.T) {
	mockRunner := NewMockRunner("")
	git := NewCustomGit(mockRunner)

	repo := &Repo{Path: "/test/"}
	git.Status(repo)

	if mockRunner.folder != "/test/" {
		t.Errorf("Folder should be /test/")
	}

	if mockRunner.command!= "git" {
		t.Errorf("Command should be git")
	}

	if len(mockRunner.args) != 2 {
		t.Errorf("Args size should be 2")
	}

	if mockRunner.args[0] != "status" {
		t.Errorf("Args 0 should be status")
	}

	if mockRunner.args[1] != "-unormal" {
		t.Errorf("Args 1 should be --all")
	}
}

func TestNotGitRepoStatus(t *testing.T) {
	output := "fatal: Not a git repository (or any of the parent directories): .git"

	mockRunner := NewMockRunner(output)
	git := NewCustomGit(mockRunner)

	repo := &Repo{Path: "/test/"}
	status := git.Status(repo)

	if status != "error" {
		t.Errorf("Should be error status")
	}
}

func TestRemoteSyncStatus(t *testing.T) {
	output := `On branch develop
Your branch is up-to-date with 'origin/develop'.
nothing to commit, working tree clean`

	mockRunner := NewMockRunner(output)
	git := NewCustomGit(mockRunner)

	repo := &Repo{Path: "/test/"}
	status := git.Status(repo)

	if !strings.HasPrefix(status, SYNC) {
		t.Errorf("Should be sync status")
	}
}

func TestRemoteBehindStatus(t *testing.T) {
	output := `On branch develop
Your branch is behind 'origin/develop' by 10 commits, and can be fast-forwarded.
  (use "git pull" to update your local branch)
nothing to commit, working tree clean`

	mockRunner := NewMockRunner(output)
	git := NewCustomGit(mockRunner)

	repo := &Repo{Path: "/test/"}
	status := git.Status(repo)

	if !strings.HasPrefix(status, BEHIND) {
		t.Errorf("Should be behind status")
	}
}

func TestRemoteAheadStatus(t *testing.T) {
	output := `On branch develop
Your branch is ahead of 'origin/develop' by 1 commit.
  (use "git push" to publish your local commits)
nothing to commit, working tree clean`

	mockRunner := NewMockRunner(output)
	git := NewCustomGit(mockRunner)

	repo := &Repo{Path: "/test/"}
	status := git.Status(repo)

	if !strings.HasPrefix(status, AHEAD) {
		t.Errorf("Should be ahead status")
	}
}

func TestRemoteDivergedStatus(t *testing.T) {
	output := `On branch chore/test
Your branch and 'origin/chore/test' have diverged,
and have 1 and 1 different commits each, respectively.
  (use "git pull" to merge the remote branch into yours)
nothing to commit, working tree clean`

	mockRunner := NewMockRunner(output)
	git := NewCustomGit(mockRunner)

	repo := &Repo{Path: "/test/"}
	status := git.Status(repo)

	if !strings.HasPrefix(status, DIVERGED) {
		t.Errorf("Should be diverged status")
	}
}

func TestRemoteNoRemoteStatus(t *testing.T) {
	output := `On branch test
nothing to commit, working tree clean`

	mockRunner := NewMockRunner(output)
	git := NewCustomGit(mockRunner)

	repo := &Repo{Path: "/test/"}
	status := git.Status(repo)

	if !strings.HasPrefix(status, NO_REMOTE) {
		t.Errorf("Should be no remote status")
	}
}

func TestLocalSyncStatus(t *testing.T) {
	output := `On branch develop
Your branch is up-to-date with 'origin/develop'.
nothing to commit, working tree clean`

	mockRunner := NewMockRunner(output)
	git := NewCustomGit(mockRunner)

	repo := &Repo{Path: "/test/"}
	status := git.Status(repo)

	if !strings.HasSuffix(status, SYNC) {
		t.Errorf("Should be sync status")
	}
}

func TestLocalUntrackedStatus(t *testing.T) {
	output := `On branch develop
Your branch is up-to-date with 'origin/develop'.
Untracked files:
  (use "git add <file>..." to include in what will be committed)

	untracked-file

nothing added to commit but untracked files present (use "git add" to track)`

	mockRunner := NewMockRunner(output)
	git := NewCustomGit(mockRunner)

	repo := &Repo{Path: "/test/"}
	status := git.Status(repo)

	if !strings.HasSuffix(status, UNTRACKED) {
		t.Errorf("Should be untracked status")
	}
}

func TestLocalChangedStatus(t *testing.T) {
	output := `On branch develop
Your branch is up-to-date with 'origin/develop'.
Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git checkout -- <file>..." to discard changes in working directory)

	modified:   file.txt

no changes added to commit (use "git add" and/or "git commit -a")`

	mockRunner := NewMockRunner(output)
	git := NewCustomGit(mockRunner)

	repo := &Repo{Path: "/test/"}
	status := git.Status(repo)

	if !strings.HasSuffix(status, CHANGED) {
		t.Errorf("Should be changed status")
	}
}

func NewMockRunner(output string) *MockRunner {
	return &MockRunner{output: output}
}

func (runner *MockRunner) Run(folder string, command string, args []string) (string, error) {
	runner.folder = folder
	runner.command = command
	runner.args = args
	return runner.output, nil
}

type MockRunner struct {
	folder string
	command string
	args []string

	output string
}

