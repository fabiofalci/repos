package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	f, err := os.Open("/home/fabio/.config/repo/repos")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var repos []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		repos = append(repos, scanner.Text())
	}

	longestName := 0
	for _, repo := range repos {
		len := len(repoName(repo))
		if len > longestName {
			longestName = len
		}
	}
	longestName = longestName + 1

	for _, repo := range repos {
		repoName := repoName(repo)
		st := status(repo)
		br := branch(repo)

		fmt.Printf("%"+strconv.Itoa(longestName)+"s %s [%s]\n", repoName, st, br)
	}
}

func repoName(repo string) string {
	return repo[strings.LastIndex(repo, "/")+1:]
}

func branch(repo string) string {
	output, err := run(repo, []string{"rev-parse", "--abbrev-ref", "HEAD"})
	if err != nil {
		fmt.Println(err)
	}
	return strings.TrimSuffix(output, "\n")
}

func status(repo string) string {
	output, err := run(repo, []string{"status", "-unormal"})
	if err != nil {
		fmt.Println(err)
	}

	if strings.Contains(output, "Not a git repo") {
		fmt.Println("Not a git repo")
		return "not a repo"
	}

	if strings.Contains(output, "have diverged") {
		return "DIVER"
	} else if strings.Contains(output, "branch is ahead") {
		return "AHEAD"
	} else if strings.Contains(output, "branch is behind") {
		return "BEHIN"
	} else if strings.Contains(output, "nothing added to commit but untracked files present") {
		return "DIRTY"
	} else if strings.Contains(output, "nothing to commit") {
		if strings.Contains(output, "is up-to-date with") {
			return "-----"
		}
		return "NO-RE"
	} else {
		return "CHANG"
	}

	return "?"
}

func run(folder string, command []string) (string, error) {
	cmd := exec.Command("git", command...)
	cmd.Dir = folder
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	if stderr.Len() > 0 {
		fmt.Println("#### Stderr ####")
		fmt.Println(stderr.String())
	}
	return out.String(), nil
}
