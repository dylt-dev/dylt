package lib

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dylt-dev/dylt/common"
)

type statusInfo struct {
	etcdDomain string
	isColima bool
	isConfigFile bool
	isIncus bool
	isVm bool
}

func checkColima () bool {
	shellPath := getShellPath()
	sCmd := fmt.Sprintf("command -v %s", "colima")
	cmd:= exec.Command(shellPath, "-c", sCmd)
	err := cmd.Run()
	var exists bool = (err == nil)

	return exists
}

func getCommandPath (cmd string) (string, error) {
	shellPath := getShellPath()
	sCmd := fmt.Sprintf("command -v %s", cmd)
	execCmd:= exec.Command(shellPath, "-c", sCmd)
	stdout, err := execCmd.StdoutPipe()
	if err != nil { return "", err }
	r := bufio.NewReader(stdout)
	err = execCmd.Start()
	if err != nil { return "", err }
	cmdPath, err := r.ReadString('\n')
	if err != nil { return "", err }
	err = execCmd.Wait()
	if err != nil { return "", err }
	cmdPath = strings.TrimSpace(cmdPath)

	return cmdPath, nil
}

func getShellPath () string {
	var shellPaths = []string {
		filepath.FromSlash("/bin/sh"),
	}

	for _, shellPath := range shellPaths {
		common.Logger.Debugf("shellPath=%s\n", shellPath)
		_, err := os.Stat(shellPath)
		if err == nil { return shellPath }
	}

	// No shell path was found
	return ""
}

func isShellAvailable () bool {
	var shellPaths = []string {
		filepath.FromSlash("/bin/sh"),
	}

	for _, shellPath := range shellPaths {
		common.Logger.Debugf("shellPath=%s\n", shellPath)
		_, err := os.Stat(shellPath)
		if err == nil { return true }
	}

	// No shell path was found
	return false
}