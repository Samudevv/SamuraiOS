package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	var full bool
	if len(os.Args) > 1 {
		if strings.ToLower(os.Args[1]) == "full" {
			full = true
		}
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	if err = os.MkdirAll(filepath.Join(homeDir, "Bilder", "Screenshots"), 0755); err != nil {
		panic(err)
	}

	fileName := time.Now().Format("screenshot-2006-01-02-15:04:05.png")

	if !full {
		var slurpOut strings.Builder

		slurp := exec.Command("slurp", "-d")
		slurp.Stdout = &slurpOut
		slurp.Stderr = os.Stderr

		if err = slurp.Run(); err != nil {
			panic(err)
		}

		geometry := strings.TrimSpace(slurpOut.String())

		grim := exec.Command("grim", "-g", geometry, filepath.Join("/tmp", fileName))
		grim.Stdout = os.Stdout
		grim.Stderr = os.Stderr

		if err = grim.Run(); err != nil {
			panic(err)
		}
	} else {
		outputName, err := exec.Command("sh", "-c", "hyprctl activeworkspace | head -n1 | cut -d ' ' -f7 | cut -d ':' -f1").Output()
		if err != nil {
			panic(err)
		}
		grim := exec.Command("grim", "-o", strings.TrimSpace(string(outputName)), filepath.Join("/tmp", fileName))
		grim.Stdout = os.Stdout
		grim.Stderr = os.Stderr

		if err = grim.Run(); err != nil {
			panic(err)
		}
	}

	swappy := exec.Command("swappy", "-o", filepath.Join(homeDir, "Bilder", "Screenshots", fileName), "-f", filepath.Join("/tmp", fileName))
	swappy.Stdout = os.Stdout
	swappy.Stderr = os.Stderr

	if err = swappy.Run(); err != nil {
		panic(err)
	}
}
