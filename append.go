package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func exeAppendFile(command, filename string) {
	words := strings.Split(command, " ")
	if len(words) == 0 {
		fmt.Fprintln(os.Stderr, "No Command")
		os.Exit(1)
	}

	var args []string
	if len(words) > 1 {
		args = words[1:]
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open \"%s\" for \"%s\": %s\n", filename, command, err)
		os.Exit(1)
	}
	defer file.Close()

	cmd := exec.Command(words[0], args...)
	cmd.Stdout = file
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Println(fmt.Sprint(command, " >> ", filename))

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "\"%s\" failed: %s\n", command, err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Invalid Arguments")
		os.Exit(1)
	}

	command := os.Args[1]
	file := os.Args[2]

	exeAppendFile(command, file)
}
