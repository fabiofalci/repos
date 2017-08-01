package main

import (
	"bytes"
	"os/exec"
)

type CommandLineRunner interface {
	Run(folder string, command string, args []string) (string, error)
}

type DefaultCommandLineRunner struct {
}

func (runner *DefaultCommandLineRunner) Run(folder string, command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
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
