package lib

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"

	"github.com/dylt-dev/dylt/common"
)

type statusInfo struct {
	etcdDomain string
	isColima bool
	isConfigFile bool
	isIncus bool
	isVm bool
}

func RunStatus () error {
	var err error
	status := new(statusInfo)

	status.isColima, err = checkColima()
	if err != nil { return err }

	status.isConfigFile, err = isExistConfigFile()
	if err != nil { return err }

	status.isIncus, err = checkIncus()
	if err != nil { return err }

	fmt.Printf("%#v\n", status)

	isShellAvailable := isShellAvailable()
	common.Logger.Debugf("isShellAvilable: %t", isShellAvailable)

	return nil
}

func checkIncus () (bool, error) {
	return isIncusAvailable()
}

func isExistConfigFile () (bool, error) {
	cfgFilePath := common.GetConfigFilePath()
	fi, err := os.Stat(cfgFilePath)
	if err != nil { return false, err }
	if fi.IsDir() {
		return false, fmt.Errorf("config file path exists, but is a directory (%s)", cfgFilePath)
	}
	if fi.Mode() & fs.ModeType > 0 {
		return false, fmt.Errorf("config file exists, but its mode is invalid (%d)", fi.Mode())
	}
	return true, nil
}

func createUnixSocketClient (socketPath string) *http.Client{
	raddr, _ := net.ResolveUnixAddr("unix", socketPath)

	// dialer := net.Dialer{}
	cli := &http.Client {
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return net.DialUnix("unix", nil, raddr)
			},
		},
	}

	return cli
}

func checkColima () (bool, error) {
	var exists = isCommandExist("colima")

	return exists, nil
}

func createShellCmd (sCmd string) *exec.Cmd {
	shellPath := getShellPath()
	shellArgs := []string{"-c", sCmd}
	cmd := exec.Command(shellPath, shellArgs...)

	return cmd
}

func isIncusAvailable () (bool, error) {
	url := "http://incus"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil { return false, err }
	
	socketPath := getIncusSocketPath()
	cli := createUnixSocketClient(socketPath)
	resp, err := cli.Do(req)
	if err != nil { return false, err }
	
	buf2, err := io.ReadAll(resp.Body)
	if err != nil { return false, err }
	fmt.Printf("buf2=%s\n", buf2)

	return true, nil
}

func isShellAvailable () bool {
	shellPaths := getShellPaths()
	for _, shellPath := range shellPaths {
		common.Logger.Debugf("shellPath=%s\n", shellPath)
		_, err := os.Stat(shellPath)
		if err == nil { return true }
	}

	// No shell path was found
	return false
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
