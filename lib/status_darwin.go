package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dylt-dev/dylt/common"
)

func createShellCmd (sCmd string) *exec.Cmd {
	shellPath := getShellPath()
	shellArgs := []string{"-c", sCmd}
	cmd := exec.Command(shellPath, shellArgs...)

	return cmd
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

func runWithOutput (cmd *exec.Cmd) (*bytes.Buffer, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil { return nil, err }
	r := bufio.NewReader(stdout)
	err = cmd.Start()
	if err != nil { return nil, err }
	var buf []byte
	buf, err = io.ReadAll(r)
	if err != nil { return nil, err }
	var buffer *bytes.Buffer
	buffer = bytes.NewBuffer(buf)
	err = cmd.Wait()
	if err != nil { return buffer, err }
	
	return buffer, nil
}

func runWithStdoutAndStderr (cmd *exec.Cmd) (*bytes.Buffer, *bytes.Buffer, error) {
	var stdout, stderr io.ReadCloser
	var err error

	stdout, err = cmd.StdoutPipe()
	if err != nil { return nil, nil, err }
	rStdout := bufio.NewReader(stdout)
	
	stderr, err = cmd.StderrPipe()
	if err != nil { return nil, nil, err }
	rStderr := bufio.NewReader(stderr)

	err = cmd.Start()
	if err != nil { return nil, nil, err }

	var bufStdout, bufStderr []byte
	var bufferStdout, bufferStderr *bytes.Buffer

	bufStdout, err = io.ReadAll(rStdout)
	if err != nil { return nil, nil, err }
	bufferStdout = bytes.NewBuffer(bufStdout)

	bufStderr, err = io.ReadAll(rStderr)
	if err != nil { return nil, nil, err }
	bufferStderr = bytes.NewBuffer(bufStderr)

	err = cmd.Wait()
	if err != nil { return bufferStdout, bufferStderr, err }
	
	return bufferStdout, bufferStderr, nil
}

func isCommandExist (cmd string) bool {
	shellPath := getShellPath()
	sCmd := fmt.Sprintf("command -v %s", cmd)
	shellCmd := exec.Command(shellPath, "-c", sCmd)
	err := shellCmd.Run()
	var exists bool = (err == nil)

	return exists
}

func getIncusSocketPath () string {
	homePath := os.Getenv("HOME")
	socketPath := filepath.Join(homePath, filepath	.FromSlash(".colima/default/incus.sock"))

	return socketPath
}

func getShellPaths () []string {
	var shellPaths = []string {
		filepath.FromSlash("/bin/sh"),
	}

	return shellPaths
}
