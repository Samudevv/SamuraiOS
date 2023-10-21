package main

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	SelectTypeRegion = iota
	SelectTypeWindow = iota
	SelectTypeOutput = iota
)

func getSmelConfig() []string {
	homeDir, _ := os.UserHomeDir()

	var smelArgs []string
	smelConf, err := os.Open(filepath.Join(homeDir, ".config/smel"))
	if err == nil {
		defer smelConf.Close()
		scanner := bufio.NewScanner(smelConf)
		for scanner.Scan() {
			smelArgs = append(smelArgs, scanner.Text())
		}
	}
	return smelArgs
}

func main() {
	var selectType int
	if len(os.Args) > 1 {
		if strings.EqualFold(os.Args[1], "full") {
			selectType = SelectTypeOutput
		} else if strings.EqualFold(os.Args[1], "window") {
			selectType = SelectTypeWindow
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

	smelArgs := []string{
		"-so",
		filepath.Join("/tmp", fileName),
	}

	if selectType == SelectTypeRegion {
		smelArgs = append(smelArgs, getSmelConfig()...)

		smel := exec.Command("smel", smelArgs...)
		smel.Stdout = os.Stdout
		smel.Stderr = os.Stderr
		if err = smel.Run(); err != nil {
			panic(err)
		}
	} else if selectType == SelectTypeOutput {
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
	} else if selectType == SelectTypeWindow {
		smelArgs = append(smelArgs, "-r")
		smelArgs = append(smelArgs, getSmelConfig()...)

		smel := exec.Command("smel", smelArgs...)
		smel.Stdout = os.Stdout
		smel.Stderr = os.Stderr
		if err = smel.Run(); err != nil {
			panic(err)
		}
	} else {
		panic("Invalid select type")
	}

	swappy := exec.Command("swappy", "-o", filepath.Join(homeDir, "Bilder", "Screenshots", fileName), "-f", filepath.Join("/tmp", fileName))
	swappy.Stdout = os.Stdout
	swappy.Stderr = os.Stderr

	if err = swappy.Run(); err != nil {
		panic(err)
	}
}
