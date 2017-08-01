package main

import (
	"os/exec"
	"bytes"
)

type Runner interface {
	Run(folder string, command string, args []string) error
}

type DefaultRunner struct {
}


func (runner *DefaultRunner) Run(folder string, command string, args []string) (string, error) {
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

