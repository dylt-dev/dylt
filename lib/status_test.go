package lib

import (
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

func TestStat0 (t *testing.T) {
	var fi os.FileInfo
	var err error

	fi, err = os.Stat("/bin/sh")
	require.NoError(t, err)
	require.NotEmpty(t, fi)

	fi, err = os.Stat("/opt/homebrew/bin/colima")
	require.NoError(t, err)
	require.NotEmpty(t, fi)

}

func TestGetColimaPath (t *testing.T) {
	common.InitLogging()
	
	colimaPath, err := getCommandPath("colima")
	require.NoError(t, err)
	fi, err := os.Stat(filepath.FromSlash(colimaPath))
	require.NotEmpty(t, fi)
	require.NoError(t, err)

	t.Logf("colimaPath=%s", colimaPath)
}

func TestColimaStatus (t *testing.T) {
	common.InitLogging()
	
	cmd := exec.Command("colima", "status")
	bufferStdout, bufferStderr, err := runWithStdoutAndStderr(cmd)
	
	t.Log(bufferStdout.String())
	t.Log(bufferStderr.String())
	
	require.NoError(t, err)
	require.NotEmpty(t, bufferStdout)
	require.NotEmpty(t, bufferStderr)
}


func TestGetUnixSocketAddreess (t *testing.T) {
	socketPath := getIncusSocketPath()
	raddr, err := net.ResolveUnixAddr("unix", socketPath)
	require.NoError(t, err)
	t.Logf("raddr=%#v", raddr)
}


func TestIsExistConfigFile (t *testing.T) {
	isExist, err := isExistConfigFile()
	require.NoError(t, err)
	require.True(t, isExist)
}

func TestStartIncus (t *testing.T) {
	cmdName := "colima"
	args := []string{"start", "--runtime", "incus"}
	cmd := exec.Command(cmdName, args...)
	stdout, stderr, err := runWithStdoutAndStderr(cmd)
	require.NoError(t, err)
	t.Log("stdout")
	t.Log(stdout.String())
	t.Log("stderr")
	t.Log(stderr.String())
}

func TestIsIncusAvailable (t *testing.T) {
	is, err := isIncusAvailable()
	require.NoError(t, err)
	require.True(t, is)
}
