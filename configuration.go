package main

import (
	"bufio"
	"log"
	"os"
	"os/user"
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
		conf.Repos = append(conf.Repos, NewRepo(scanner.Text()))
	}
}

func (conf *Configuration) GetLongestName() int {
	longestName := 0
	for _, repo := range conf.Repos {
		if len(repo.Path) > 0 && string(repo.Path[0]) != "#" {
			len := len(repo.Name())
			if len > longestName {
				longestName = len
			}
		}
	}
	return longestName + 1
}
