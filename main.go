package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"runtime"
	"strconv"
	"time"
)

var version string
var commit string
var buildDate string

const (
	UNTRACKED = "UNTRA"
	SYNC      = "-----"
	CHANGED   = "CHANG"
	DIVERGED  = "DIVER"
	AHEAD     = "AHEAD"
	BEHIND    = "BEHIN"
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
	var clone bool = false
	var branches bool = false
	var hideUnchanged bool = false
	var profile string = ""

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "clone, c",
			Usage:       "Execute git clone <url> <dir>",
			Destination: &clone,
		},
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
		cli.StringFlag{
			Name:        "profile, p",
			Usage:       "Show profile repos",
			Destination: &profile,
		},
		cli.BoolFlag{
			Name:        "hide, hi",
			Usage:       "Hide repos if on master/developer and unchanged",
			Destination: &hideUnchanged,
		},
	}

	app.Action = func(c *cli.Context) error {
		if clone {
			cloneRepos(profile)
		} else {
			showRepos(profile, fetch, branches, hideUnchanged)
		}
		return nil
	}

	app.Run(os.Args)
}

func cloneRepos(profile string) {
	conf := NewConfiguration(profile)
	git := NewDefaultGit()

	longestName := conf.GetLongestName()

	fmt.Printf("Cloning repos...\n")
	for _, repo := range conf.Repos {
		repoName := repo.Name()
		cloned := git.Clone(repo)
		if cloned != "error" {
			fmt.Printf("%"+strconv.Itoa(longestName)+"s cloned\n", repoName)
		}
	}
	fmt.Println("Done")
}

func showRepos(profile string, fetchRepo bool, showBranches bool, hideUnchanged bool) {
	conf := NewConfiguration(profile)
	git := NewDefaultGit()

	longestName := conf.GetLongestName()

	fmt.Printf("%"+strconv.Itoa(longestName)+"s Remot Local [branch]\n", "")
	for _, repo := range conf.Repos {
		if fetchRepo {
			git.Fetch(repo)
		}
		repoName := repo.Name()
		st := git.Status(repo)
		if st != "error" {
			br := git.Branch(repo)
			if hideUnchanged && st == "----- -----" && (br == "develop" || br == "master") {
				continue
			}
			if showBranches {
				brs := git.Branches(repo, br)
				fmt.Printf("%"+strconv.Itoa(longestName)+"s %s [%s] %s\n", repoName, st, br, brs)
			} else {
				fmt.Printf("%"+strconv.Itoa(longestName)+"s %s [%s]\n", repoName, st, br)
			}
		} else {
			fmt.Printf("%"+strconv.Itoa(longestName)+"s error\n", repoName)
		}
	}
}
