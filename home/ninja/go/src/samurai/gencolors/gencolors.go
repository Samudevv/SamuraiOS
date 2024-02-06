package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	colorFile := os.Args[1]
	genType := strings.ToLower(os.Args[2])

	colors, err := ParseColors(colorFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch genType {
	case "mako":
		configTmpl := os.Args[3]
		config := os.Args[4]

		configTmplFile, err := os.Open(configTmpl)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer configTmplFile.Close()

		configFile, err := os.Create(config)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer configFile.Close()

		if err = GenerateMako(colors, configTmplFile, configFile); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "Invalid genType")
		os.Exit(1)
	}

	fmt.Println("Sucessfully generated", genType)
}
