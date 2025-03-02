package lib

import (
	"log"
	"log/slog"
	"os"
	"os/exec"
)

const PATH_DaylightScript = "/opt/bin/daylight.sh"


func CallDaylightScript (args []string) (int, error) {
	_, err := os.Stat(PATH_DaylightScript)
	if err != nil { return 1, err }
	cmd := exec.Command(PATH_DaylightScript, args...)
	err = cmd.Run()
	var rc int
	if err != nil {
		log.Fatalf("rc=%d\n", err.(*exec.ExitError).ExitCode())
		rc = err.(*exec.ExitError).ExitCode()
	} else {
		rc = 0
	}
	return rc, nil
}


func CallDaylightScriptO (args []string) (int, []byte, error) {
	_, err := os.Stat(PATH_DaylightScript)
	if err != nil { return 1, []byte{}, err }
	cmd := exec.Command(PATH_DaylightScript, args...)
	stdout, err := cmd.Output()
	var rc int
	if err != nil {
		log.Fatalf("rc=%d\n", err.(*exec.ExitError).ExitCode())
		stdout = []byte{}
		rc = err.(*exec.ExitError).ExitCode()
	} else {
		rc = 0
	}
	return rc, stdout, nil
}


func IsPathExecutable (path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil { return false, err }
	perm := fi.Mode().Perm()
	if perm & 0111 > 0 {
		return true, nil
	}
	return false, nil	
}

func RunCommand (cmdName string, args... string) (int, []byte, error) {
	slog.Debug("lib.RunCommand()", "cmdName", cmdName, "args", args)
	// execute command with args + read stdout into a string
	cmd := exec.Command(cmdName, args...)
	s, err := cmd.Output()
	var rc int
	// on error, set rc to non-zero and stdout to empty
	if err != nil {
		switch err := err.(type) {
		case *exec.Error: rc = 1
		case *exec.ExitError: rc = err.ExitCode()
		default: rc = 1
		}
		slog.Error("command failed", "rc", rc)
		s = []byte{}
	} else {
		rc = 0
		err = nil
	}

	return rc, s, err
}

func RunScript (path string, args []string) (int, []byte, error) {
	// Check if script exists
	flag, err := IsPathExecutable(path)
	if !flag || err != nil {
		if err != nil {
			slog.Error(err.Error())
		}
		return 1, []byte{}, err
	}
	// delegate heavy lifting
	return RunCommand(path, args...)
}