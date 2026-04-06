package lib

import (
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	RunAndTestCommand(t, "dylt status")
}

// Simple test to test `os.Stat()` which we hadn't used before.
func TestStat0(t *testing.T) {
	var fi os.FileInfo
	var err error

	fi, err = os.Stat("/bin/sh")
	require.NoError(t, err)
	require.NotEmpty(t, fi)

	fi, err = os.Stat("/opt/homebrew/bin/colima")
	require.NoError(t, err)
	require.NotEmpty(t, fi)

}

// Check that `colima` is on the PATH and that its PATH entry exists
// @note `os/exec.LookPath()` might be better
func TestGetColimaPath(t *testing.T) {
	common.InitLogging()

	colimaPath, err := getCommandPath("colima")
	require.NoError(t, err)
	fi, err := os.Stat(filepath.FromSlash(colimaPath))
	require.NotEmpty(t, fi)
	require.NoError(t, err)
	isExecutable := fi.Mode().Perm()&0x111 > 0
	require.True(t, isExecutable)
	t.Logf("colimaPath=%s", colimaPath)
}

// Get the incus socket path from config & validate that it exists
func TestGetUnixSocketAddreess(t *testing.T) {
	socketPath := getIncusSocketPath()
	raddr, err := net.ResolveUnixAddr("unix", socketPath)
	require.NoError(t, err)
	t.Logf("raddr=%#v", raddr)
}

// Test that the dylt config file exists
func TestIsExistConfigFile(t *testing.T) {
	isExist, err := isExistConfigFile()
	require.NoError(t, err)
	require.True(t, isExist)
}

// Test if incus is available
func TestIncusIsAvailable(t *testing.T) {
	is, err := isIncusAvailable()
	require.NoError(t, err)
	require.True(t, is)
}

// Test if the incus `dylt` container exists
// @note I'm not actually sure what the `dylt` container is for
func TestIncusIsDyltContainerExist(t *testing.T) {
	if os.Getenv("DYLT_TEST_SYSTEST") == "" {
		t.Skip("sys test only")
	}
	
	flag, err := isIncusDyltContainerExist()
	require.NoError(t, err)
	require.True(t, flag)
}

func TestIncusListContainerNames(t *testing.T) {
	names, err := getIncusContainerNames()
	require.NoError(t, err)
	require.NotEmpty(t, names)
	for _, name := range names {
		t.Log(name)
	}
}