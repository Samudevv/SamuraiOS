package main

import (
	"fmt"
	"os"
)

func main() {
	fileName := os.Args[1]

	colors, err := ParseColors(fileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = GenerateSCSS(colors, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
