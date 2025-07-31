package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dylt-dev/dylt/color"
	"github.com/dylt-dev/dylt/common"
)

func RunStatus () error {
	fmt.Println("hi")
	var err error
	status := new(statusInfo)

	status.isConfigFile, err = isExistConfigFile()
	if err != nil {
		common.Logger.Debug(err.Error())
		status.isConfigFile = false
	}

	status.isIncusActive, err = isIncusActive()
	if err != nil {
		common.Logger.Debug(err.Error())
		status.isIncusActive = false
	}

	status.isVm, err = isIncusDyltContainerExist()
	if err != nil {
		common.Logger.Debug(err.Error())
		status.isVm = false
	}

	fmt.Printf("%-42s %s\n", string(common.Highlight("is config file exist")), color.StyleBool(status.isConfigFile))
	fmt.Printf("%-42s %s\n", string(common.Highlight("is incus active")), color.StyleBool(status.isIncusActive))
	fmt.Printf("%-42s %s\n", string(common.Highlight("is incus dylt container exists")), color.StyleBool(status.isVm))
	
	return nil
}



// Run /bin/sh -c 'command -v $cmd' to get the path of a command
// Return an error if cmd is not installed
func getCommandPath (cmdName string) (string, error) {
	sCmd := fmt.Sprintf("command -v %s", cmdName)
	shellCmd := createShellCmd(sCmd)
	buffer, err := runWithOutput(shellCmd)
	if err != nil { return "", err }
	
	cmdPath, err := buffer.ReadString('\n')
	if err != nil { return "", err }
	cmdPath = strings.TrimSpace(cmdPath)

	return cmdPath, nil
}

// Iterate over a collection of shell paths (eg /bin/sh), returning the first
// path found on the host.
//
// return "" if no shell path was found
func getShellPath () string {
	// All potential shell paths
	var shellPaths = getShellPaths()

	for _, shellPath := range shellPaths {
		common.Logger.Debugf("shellPath=%s\n", shellPath)
		_, err := os.Stat(shellPath)
		if err == nil { return shellPath }
	}

	// No shell path was found
	return ""
}

func getShellPaths () []string {
	var shellPaths = []string {
		filepath.FromSlash("/bin/sh"),
	}

	return shellPaths
}

func getIncusSocketPath () string {
	socketPath := "/var/lib/incus/unix.socket"

	return socketPath
}

// func isCommandExist (cmd string) bool {
// 	shellPath := getShellPath()
// 	sCmd := fmt.Sprintf("command -v %s", cmd)
// 	shellCmd := exec.Command(shellPath, "-c", sCmd)
// 	err := shellCmd.Run()
// 	var exists bool = (err == nil)

// 	return exists
// }
