package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
	wfPID, err := getPID("wf-recorder")
	if err != nil {
		notifyError("wf-recorder: ", err)
		os.Exit(1)
	}

	if wfPID != 0 {
		var killErr strings.Builder
		kill := exec.Command("kill", "-s", "SIGINT", strconv.Itoa(wfPID))
		kill.Stderr = &killErr
		if err := kill.Run(); err != nil {
			notifyError("Failed to execute kill: ", killErr.String())
			os.Exit(1)
		}
		notify("Recording Stopped")
	} else {
		smelArgs := getSmelConfig()

		var smelOut, smelErr strings.Builder
		smel := exec.Command("smel", smelArgs...)
		smel.Stdout = &smelOut
		smel.Stderr = &smelErr
		if err := smel.Run(); err != nil {
			smelErrStr := strings.TrimSpace(smelErr.String())
			if smelErrStr == "selection cancelled" {
				notify("Recording Cancelled")
				os.Exit(0)
			}
			notifyError("Failed to execute slurp: ", smelErrStr)
			os.Exit(1)
		}

		geo := strings.TrimSpace(smelOut.String())

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
				notifyError("Failed to execute wf-recorder: ", strings.TrimSpace(wfErr.String()))
				os.Exit(1)
			}

			notify("Converting to GIF ...")
			var gifErr strings.Builder
			gifski := exec.Command("gifski", "-o", "recording.gif", "recording.mkv")
			gifski.Stderr = &gifErr
			if err := gifski.Run(); err != nil {
				notifyError("Failed to execute gifski: ", strings.TrimSpace(gifErr.String()))
				os.Exit(1)
			}
			os.Remove("recording.mkv")

			filename := time.Now().Format("recording-2006-01-02-15:04:05.gif")
			homeDir, _ := os.UserHomeDir()
			foldername := filepath.Join(homeDir, "Videos")
			os.Mkdir(foldername, 0755)
			newPath := filepath.Join(foldername, filename)
			mv := exec.Command("mv", "recording.gif", newPath)
			mv.Run()
			notify("Recording saved to ", newPath)

			eog := exec.Command("eog", newPath)
			eog.Run()
		} else {
			notify("Recording Cancelled")
		}
	}
}

// Returns 0 if no PID can be found
func getPID(name string) (int, error) {
	ps := exec.Command("ps", "-e")
	grep := exec.Command("grep", name)

	var psOut, grepOut strings.Builder
	var errOut strings.Builder

	ps.Stdout = &psOut
	ps.Stderr = &errOut
	grep.Stdout = &grepOut
	grep.Stderr = &errOut

	if err := ps.Run(); err != nil {
		return 0, fmt.Errorf("Failed to execute ps: %s", errOut.String())
	}

	grepIn := strings.NewReader(psOut.String())
	grep.Stdin = grepIn

	if err := grep.Run(); err != nil {
		return 0, nil
	}

	grepStr := strings.TrimSpace(grepOut.String())
	grepWords := strings.Split(grepStr, " ")
	pidStr := grepWords[0]

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
