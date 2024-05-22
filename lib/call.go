package lib

import (
	"log"
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