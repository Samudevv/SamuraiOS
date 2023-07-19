package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	filename := os.Args[1]
	old := os.Args[2]
	new := os.Args[3]

	filestat, err := os.Stat(filename)
	if err != nil {
		logError(err)
		os.Exit(1)
	}

	input, err := ioutil.ReadFile(filename)
	if err != nil {
		logError(err)
		os.Exit(1)
	}

	output := bytes.Replace(input, []byte(old), []byte(new), 1)

	err = ioutil.WriteFile(filename, output, filestat.Mode())
	if err != nil {
		logError(err)
		os.Exit(1)
	}
}

func logError(msg ...any) {
	msgStr := fmt.Sprint(msg...)
	fmt.Fprintf(os.Stderr, "\n\n\033[30;41m[ERROR]\033[0;33m %s\033[0m\n\n", msgStr)
}
