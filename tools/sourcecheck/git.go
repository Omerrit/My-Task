package main

import (
	"bytes"
	"errors"
	"os/exec"
)

func findGitRepos(dir string) (string, error) {
	reposDir, err := exec.Command("git", "-C", dir, "rev-parse", "--show-toplevel").Output()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return "", nil
		}
		return "", err
	}
	reposDir = bytes.Trim(reposDir, "\n")
	return string(reposDir), nil
}
