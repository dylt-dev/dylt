package lib

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dylt-dev/dylt/color"
	"github.com/dylt-dev/dylt/common"
)

func RunStatus() error {
	var err error
	status := new(statusInfo)

	status.isColimaExist, err = checkColima()
	fmt.Printf("%-20s %-5t %v\n", "isColimaExist", status.isColimaExist, err)
	if err != nil {
		return err
	}

	status.isColimaActive, err = isColimaActive()
	fmt.Printf("%-20s %-5t %v\n", "isColimaActive", status.isColimaActive, err)
	if err != nil {
		return err
	}

	status.isConfigFile, err = isExistConfigFile()
	fmt.Printf("%-20s %-5t %v\n", "isConfigFile", status.isConfigFile, err)
	if err != nil {
		return err
	}

	status.isIncusActive, err = isIncusActive()
	fmt.Printf("%-20s %-5t %v\n", "isIncusActive", status.isIncusActive, err)
	if err != nil {
		return err
	}

	status.isVm, err = isIncusDyltContainerExist()
	fmt.Printf("%-20s %-5t %v\n", "IsIncusDyltContainerExist", status.isVm, err)
	if err != nil {
		return err
	}

	fmt.Printf("%#v\n", status)

	isShellAvailable := isShellAvailable()
	common.Logger.Debugf("isShellAvilable: %t", isShellAvailable)

	fmt.Println()
	fmt.Printf("%-42s %s\n", string(common.Highlight("is colima exist")), styleBool(status.isColimaExist))
	fmt.Printf("%-42s %s\n", string(common.Highlight("is colima available")), styleBool(status.isColimaActive))
	fmt.Printf("%-42s %s\n", string(common.Highlight("is config file exist")), styleBool(status.isConfigFile))
	fmt.Printf("%-42s %s\n", string(common.Highlight("is incus active")), styleBool(status.isIncusActive))
	fmt.Printf("%-42s %s\n", string(common.Highlight("is incus dylt container exists")), styleBool(status.isVm))

	return nil
}

func getColimaSocketPath() string {
	homePath := os.Getenv("HOME")
	socketRelPath := filepath.FromSlash(".colima/_lima/colima/ssh.sock")
	socketPath := filepath.Join(homePath, socketRelPath)

	return socketPath
}

// Run /bin/sh -c 'command -v $cmd' to get the path of a command
// Return an error if cmd is not installed
func getCommandPath(cmdName string) (string, error) {
	sCmd := fmt.Sprintf("command -v %s", cmdName)
	shellCmd := createShellCmd(sCmd)
	buffer, err := runWithOutput(shellCmd)
	if err != nil {
		return "", err
	}

	cmdPath, err := buffer.ReadString('\n')
	if err != nil {
		return "", err
	}
	cmdPath = strings.TrimSpace(cmdPath)

	return cmdPath, nil
}

func getIncusSocketPath() string {
	homePath := os.Getenv("HOME")
	socketPath := filepath.Join(homePath, filepath.FromSlash(".colima/default/incus.sock"))

	return socketPath
}

// Iterate over a collection of shell paths (eg /bin/sh), returning the first
// path found on the host.
//
// return "" if no shell path was found
func getShellPath() string {
	// All potential shell paths
	var shellPaths = getShellPaths()

	for _, shellPath := range shellPaths {
		common.Logger.Debugf("shellPath=%s\n", shellPath)
		_, err := os.Stat(shellPath)
		if err == nil {
			return shellPath
		}
	}

	// No shell path was found
	return ""
}

func getShellPaths() []string {
	var shellPaths = []string{
		filepath.FromSlash("/bin/sh"),
	}

	return shellPaths
}

func isColimaActive() (bool, error) {
	path := getColimaSocketPath()
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	mode := fi.Mode() & fs.ModeType
	if mode != fs.ModeSocket {
		err = fmt.Errorf("%s exists but is not a Unix socket", path)
		return false, err
	}
	return true, nil
}

func isCommandExist(cmd string) bool {
	shellPath := getShellPath()
	sCmd := fmt.Sprintf("command -v %s", cmd)
	shellCmd := exec.Command(shellPath, "-c", sCmd)
	err := shellCmd.Run()
	var exists bool = (err == nil)

	return exists
}

func styleBool (flag bool) color.Styledstring {
	s := color.Styledstring(strconv.FormatBool(flag))
	if flag {
		s = s.Fg(color.X11.Green)
	} else {
		s = s.Fg(color.X11.Red)
	}

	return s
}
