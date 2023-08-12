package main

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	args := parseConfig()
	if len(os.Args) > 1 {
		args = append(args, os.Args[1:]...)
	}

	slurp := exec.Command("slurp", args...)
	slurp.Stdout = os.Stdout
	slurp.Stdin = os.Stdin
	slurp.Stderr = os.Stderr

	if err := slurp.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}

func parseConfig() []string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return []string{}
	}

	file, err := os.Open(filepath.Join(homeDir, ".config", "slurp"))
	if err != nil {
		return []string{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var args []string
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.FieldsFunc(line, func(c rune) bool { return c == ' ' })
		args = append(args, words...)
	}

	return args
}
