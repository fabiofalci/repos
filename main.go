package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var version string
var commit string
var buildDate string

const (
	UNTRACKED = "UNTRA"
	SYNC = "-----"
	CHANGED = "CHANG"
	DIVERGED = "DIVER"
	AHEAD = "AHEAD"
	BEHIND = "BEHIN"
	NO_REMOTE = "NO-RE"
)

func main() {
	app := cli.NewApp()
	app.Name = "repos"
	app.Version = version
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println("Version: " + c.App.Version)
		fmt.Println("Git commit: " + commit)
		if i, err := strconv.ParseInt(buildDate, 10, 64); err == nil {
			fmt.Println("Build date: " + time.Unix(i, 0).UTC().String())
		}
		fmt.Println("Go version: " + runtime.Version())
	}

	var fetch bool = false
	var branches bool = false

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "fetch, f",
			Usage:       "Execute git fetch --all",
			Destination: &fetch,
		},
		cli.BoolFlag{
			Name:        "branches, b",
			Usage:       "Show all branches",
			Destination: &branches,
		},
	}

	app.Action = func(c *cli.Context) error {
		show(fetch, branches)
		return nil
	}

	app.Run(os.Args)
}

func show(fetchRepo bool, showBranches bool) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
		return
	}

	f, err := os.Open(usr.HomeDir + "/.config/repos/repos")
	if err != nil {
		log.Fatal("Cannot open '~/.config/repos/repos'. Have you created a repos configuration file?")
		return
	}
	defer f.Close()

	var repos []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		repos = append(repos, scanner.Text())
	}

	profiles := getProfiles(repos)

	for name, r := range profiles {
		showProfile(name, r, fetchRepo, showBranches)
	}
}

func showProfile(name string, repos []string, fetchRepo bool, showBranches bool) {
	longestName := 0
	for _, repo := range repos {
		if len(repo) > 0 && string(repo[0]) != "#" {
			len := len(repoName(repo))
			if len > longestName {
				longestName = len
			}
		}
	}
	longestName = longestName + 1

	if name != "__default__" {
		fmt.Printf("%v\n", name)
	}
	fmt.Printf("%"+strconv.Itoa(longestName)+"s Remot Local [branch]\n", "")
	for _, repo := range repos {
		if len(repo) == 0 || string(repo[0]) == "#" {
			continue
		}
		if fetchRepo {
			fetch(repo)
		}
		repoName := repoName(repo)
		st := status(repo)
		if st != "error" {
			br := branch(repo)
			if showBranches {
				brs := branches(repo, br)
				fmt.Printf("%"+strconv.Itoa(longestName)+"s %s [%s] %s\n", repoName, st, br, brs)
			} else {
				fmt.Printf("%"+strconv.Itoa(longestName)+"s %s [%s]\n", repoName, st, br)
			}
		} else {
			fmt.Printf("%"+strconv.Itoa(longestName)+"s error\n", repoName)
		}
	}
}

func getProfiles(repos []string) map[string][]string {
	profiles := make(map[string][]string)

	var current string
	if string(repos[0][0]) != "#" {
		current = "__default__"
		profiles[current] = make([]string, 0)
	}

	for _, repo := range repos {
		if len(repo) == 0 {
			continue
		}

		if string(repo[0]) == "#" {
			current = repo[1:]
			profiles[current] = make([]string, 0)
		}
		profiles[current] = append(profiles[current], repo)
	}

	return profiles
}

func fetch(repo string) {
	run(repo, []string{"fetch", "--all"})
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

func branches(repo string, currentBranch string) string {
	output, err := run(repo, []string{"branch", "--column", "--format=%(refname:short)~"})
	if err != nil {
		fmt.Println(err)
	}
	output = strings.Replace(output, "\n", "", -1)
	output = strings.Replace(output, currentBranch+"~", "", -1)
	output = strings.Replace(output, " ", "", -1)
	output = strings.Replace(output, "~", " ", -1)
	return output
}

func status(repo string) string {
	output, err := run(repo, []string{"status", "-unormal"})
	if err != nil || strings.Contains(output, "Not a git repo") {
		return "error"
	}
	return remoteStatus(output) + " " + localStatus(output)
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
	return out.String(), nil
}
