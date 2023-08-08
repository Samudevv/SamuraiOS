package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	screenCapturePID, err := getScreenCapturePID()
	if err != nil {
		notifyError(err)
		os.Exit(1)
	}

	if screenCapturePID != 0 {
		var killErr strings.Builder
		kill := exec.Command("kill", "-s", "SIGINT", strconv.Itoa(screenCapturePID))
		kill.Stderr = &killErr
		if err := kill.Run(); err != nil {
			notifyError("Failed to execute kill: ", killErr.String())
			os.Exit(1)
		}
		notify("Recording Stopped")
	} else {
		var slurpOut, slurpErr strings.Builder
		slurp := exec.Command("slurp")
		slurp.Stdout = &slurpOut
		slurp.Stderr = &slurpErr
		if err := slurp.Run(); err != nil {
			notifyError("Failed to execute slurp: ", slurpErr.String())
			os.Exit(1)
		}

		geo := strings.TrimSpace(slurpOut.String())

		if len(geo) != 0 {
			geoWords := strings.Split(geo, " ")
			geoPos := geoWords[0]
			geoSize := geoWords[1]

			notifyTimer("3")
			time.Sleep(1 * time.Second)
			notifyTimer("2")
			time.Sleep(1 * time.Second)
			notifyTimer("1")
			time.Sleep(1 * time.Second)
			notifyTimer("Start Recording")

			os.Chdir("/tmp")
			os.Remove("recording.mp4")
			var wfErr strings.Builder
			wf_recorder := exec.Command("wf-recorder", "-g", strings.Join([]string{geoPos, geoSize}, " "))
			wf_recorder.Stderr = &wfErr
			if err := wf_recorder.Run(); err != nil {
				notifyError("Failed to execute wf-recorder: ", wfErr.String())
				os.Exit(1)
			}

			notify("Converting to GIF ...")
			var gifErr strings.Builder
			gifski := exec.Command("gifski", "-o", "recording.gif", "recording.mp4")
			gifski.Stderr = &gifErr
			if err := gifski.Run(); err != nil {
				notifyError("Failed to execute gifski: ", gifErr.String())
				os.Exit(1)
			}
			os.Remove("recording.mp4")

			notify("Recording Done")

			eog := exec.Command("eog", "recording.gif")
			eog.Run()
		} else {
			notify("Recording Cancelled")
		}
	}
}

// Returns 0 if no PID can be found
func getScreenCapturePID() (int, error) {
	ps := exec.Command("ps", "-e")
	grep := exec.Command("grep", "wf-recorder")
	cut := exec.Command("cut", "-d", " ", "-f1")

	var psOut, grepOut, cutOut strings.Builder
	var errOut strings.Builder

	ps.Stdout = &psOut
	ps.Stderr = &errOut
	grep.Stdout = &grepOut
	grep.Stderr = &errOut
	cut.Stdout = &cutOut
	cut.Stderr = &errOut

	if err := ps.Run(); err != nil {
		return 0, fmt.Errorf("Failed to execute ps: %s", errOut.String())
	}

	grepIn := strings.NewReader(psOut.String())
	grep.Stdin = grepIn

	if err := grep.Run(); err != nil {
		return 0, nil
	}

	cutIn := strings.NewReader(grepOut.String())
	cut.Stdin = cutIn

	if err := cut.Run(); err != nil {
		return 0, fmt.Errorf("Failed to execute cut: %s", errOut.String())
	}

	pidStr := strings.TrimSpace(cutOut.String())

	pid, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Failed to parse PID as int: %s", err)
	}

	return int(pid), nil
}

func notifyTimer(value string) {
	fmt.Println(value)
	notify := exec.Command("notify-send", "-t", "1000", value)
	notify.Run()
}

func notify(msg ...any) {
	msgStr := fmt.Sprint(msg...)
	fmt.Println(msgStr)
	notify := exec.Command("notify-send", msgStr)
	notify.Run()
}

func notifyError(msg ...any) {
	msgStr := fmt.Sprint(msg...)
	fmt.Fprintln(os.Stderr, msgStr)
	notify := exec.Command("notify-send", "-u", "critical", msgStr)
	notify.Run()
}
