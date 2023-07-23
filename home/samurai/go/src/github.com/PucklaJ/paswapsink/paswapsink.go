package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Sink struct {
	Number  string
	Name    string
	Running bool
}

const (
	LOW      = iota
	NORMAL   = iota
	CRITICAL = iota
)

func sendNotification(title string, urgency int, message ...any) error {
	var urg string
	switch urgency {
	case LOW:
		urg = "low"
	case CRITICAL:
		urg = "critical"
	default:
		urg = "normal"
	}

	msg := fmt.Sprint(message...)

	cmd := exec.Command(
		"notify-send",
		"-u",
		urg,
		"-t",
		"5000",
		title,
		msg,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func main() {
	infoCMD := exec.Command("pactl", "list")

	var infoBuilder strings.Builder
	infoCMD.Stderr = os.Stderr
	infoCMD.Stdout = &infoBuilder

	err := infoCMD.Run()
	if err != nil {
		sendNotification("paswapsink", CRITICAL, "Failed to run info command: ", err)
		os.Exit(1)
	}

	infoStr := infoBuilder.String()

	// Parse for sinks

	var sinks []Sink
	infoScanner := bufio.NewScanner(strings.NewReader(infoStr))

	for infoScanner.Scan() {
		infoLine := infoScanner.Text()
		if strings.HasPrefix(infoLine, "Sink #") {
			number := strings.TrimPrefix(infoLine, "Sink #")
			infoScanner.Scan()
			state := strings.TrimPrefix(infoScanner.Text(), "\tState: ")
			infoScanner.Scan()
			name := strings.TrimPrefix(infoScanner.Text(), "\tName: ")

			running := state == "RUNNING"

			sink := Sink{
				Number:  number,
				Running: running,
				Name:    name,
			}

			sinks = append(sinks, sink)
		}
	}

	if len(sinks) == 0 {
		sendNotification("paswapsink", CRITICAL, "No sinks detected")
		os.Exit(1)
	}

	if len(sinks) == 1 {
		sendNotification("paswapsink", CRITICAL, "Only one sink has been found")
		os.Exit(1)
	}

	// get name of default sink
	defaultSinkCMD := exec.Command("pactl", "get-default-sink")
	var defaultSinkBuilder strings.Builder
	var defaultSinkErrBuilder strings.Builder
	defaultSinkCMD.Stdout = &defaultSinkBuilder
	defaultSinkCMD.Stderr = &defaultSinkErrBuilder
	err = defaultSinkCMD.Run()
	if err != nil {
		sendNotification("paswapsink", CRITICAL, "Failed to get default sink: ", err, "; ", defaultSinkErrBuilder.String())
		os.Exit(1)
	}

	defaultSinkName := strings.TrimSpace(defaultSinkBuilder.String())

	for i, s := range sinks {
		if s.Name == defaultSinkName {
			index := i
			if index == len(sinks)-1 {
				index = 0
			} else {
				index++
			}

			var errBuilder strings.Builder

			fmt.Println("Swapping from sink\033[1m\033[32m", s.Name, "\033[0mto\033[1m\033[32m", sinks[index].Name, "\033[0m")

			sinkCMD := exec.Command("pactl", "set-default-sink", sinks[index].Number)
			sinkCMD.Stderr = &errBuilder
			sinkCMD.Stdout = &errBuilder
			err = sinkCMD.Run()
			if err != nil {
				sendNotification("paswapsink", CRITICAL, "Failed to swap sink: ", err, "; ", errBuilder.String())
				os.Exit(1)
			}

			sendNotification("Swapped Audio Sink", LOW, sinks[index].Name)
			return
		}
	}

	sendNotification("paswapsink", CRITICAL, "Default sink has not been found")
	os.Exit(1)
}
