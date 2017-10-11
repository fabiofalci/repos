package main

import (
	"bufio"
	"log"
	"os"
	"os/user"
	"strings"
)

type Configuration struct {
	Repos []*Repo
}

func NewConfiguration() *Configuration {
	conf := &Configuration{
		Repos: make([]*Repo, 0),
	}
	conf.readFile()
	return conf
}

func (conf *Configuration) readFile() {
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

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, " ") {
			words := strings.Fields(line)
			conf.Repos = append(conf.Repos, &Repo{
				Path: words[0],
				RepoUrl: words[1],
			})
		} else {
			conf.Repos = append(conf.Repos, &Repo{Path: line})
		}
	}
}

func (conf *Configuration) GetLongestName() int {
	longestName := 0
	for _, repo := range conf.Repos {
		if len(repo.Path) > 0 && string(repo.Path[0]) != "#" {
			length := len(repo.Name())
			if length > longestName {
				longestName = length
			}
		}
	}
	return longestName + 1
}
