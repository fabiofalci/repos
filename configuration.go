package main

import (
	"bufio"
	"log"
	"os"
	"os/user"
	"strings"
)

type Configuration struct {
	Repos   []*Repo
	profile string
}

func NewConfiguration(profile string) *Configuration {
	conf := &Configuration{
		Repos: make([]*Repo, 0),
		profile: profile,
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

	configurationPath := conf.getProfileConfigurationPath()
	f, err := os.Open(usr.HomeDir + configurationPath)
	if err != nil {
		log.Fatalf("Cannot open '%v'. Have you created a repos configuration file?", configurationPath)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
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

func (conf *Configuration) getProfileConfigurationPath() string {
	if conf.profile == "" {
		return "/.config/repos/repos"
	}
	return "/.config/repos/repos." + conf.profile
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
