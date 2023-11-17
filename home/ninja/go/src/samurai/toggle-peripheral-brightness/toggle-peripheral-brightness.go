package main

import (
	"os"
	"os/exec"
	"strings"
)

const (
	modeOn     = 0
	modeOff    = 1
	modeToggle = 2
)

func main() {
	var mode int = modeToggle

	if len(os.Args) > 1 {
		switch strings.ToLower(os.Args[1]) {
		case "on":
			mode = modeOn
		case "off":
			mode = modeOff
		}
	}

	if mode == modeToggle {
		// Toggle eruption
		for _, device_id := range []string{"00", "01"} {
			eruption := exec.Command("eruptionctl", "devices", "brightness", device_id)
			var stdout strings.Builder
			eruption.Stdout = &stdout

			if err := eruption.Run(); err != nil {
				continue
			}

			lines := strings.Split(stdout.String(), "\n")
			words := strings.Split(lines[len(lines)-2], " ")
			brightness := strings.TrimSuffix(words[len(words)-1], "%")

			if brightness == "0" {
				brightness = "100"
			} else {
				brightness = "0"
			}

			eruption = exec.Command("eruptionctl", "devices", "brightness", device_id, brightness)
			eruption.Stdout = os.Stdout
			eruption.Stderr = os.Stderr
			eruption.Run()
		}

		// OpenRGB brightness
		if _, err := os.Stat("/tmp/openrgb-brightness"); err != nil && os.IsNotExist(err) {
			file, err := os.Create("/tmp/openrgb-brightness")
			if err == nil {
				file.Close()
				openrgb := exec.Command("openrgb", "--noautoconnect", "-b", "0")
				openrgb.Stdout = os.Stdout
				openrgb.Stderr = os.Stderr
				openrgb.Run()
			}
		} else if err == nil {
			if err = os.Remove("/tmp/openrgb-brightness"); err == nil {
				openrgb := exec.Command("openrgb", "--noautoconnect", "-p", "samurai")
				openrgb.Stdout = os.Stdout
				openrgb.Stderr = os.Stderr
				openrgb.Run()
			}
		}
	} else if mode == modeOn {
		for _, device_id := range []string{"00", "01"} {
			eruption := exec.Command("eruptionctl", "devices", "brightness", device_id, "100")
			eruption.Stdout = os.Stdout
			eruption.Stderr = os.Stderr
			eruption.Run()
		}

		if err := os.Remove("/tmp/openrgb-brightness"); err == nil {
			openrgb := exec.Command("openrgb", "--noautoconnect", "-p", "samurai")
			openrgb.Stdout = os.Stdout
			openrgb.Stderr = os.Stderr
			openrgb.Run()
		}
	} else if mode == modeOff {
		for _, device_id := range []string{"00", "01"} {
			eruption := exec.Command("eruptionctl", "devices", "brightness", device_id, "0")
			eruption.Stdout = os.Stdout
			eruption.Stderr = os.Stderr
			eruption.Run()
		}

		file, err := os.Create("/tmp/openrgb-brightness")
		if err == nil {
			file.Close()
			openrgb := exec.Command("openrgb", "--noautoconnect", "-b", "0")
			openrgb.Stdout = os.Stdout
			openrgb.Stderr = os.Stderr
			openrgb.Run()
		}
	}
}
