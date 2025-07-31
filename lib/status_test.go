package lib

import (
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

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

func TestGetColimaPath(t *testing.T) {
	common.InitLogging()

	colimaPath, err := getCommandPath("colima")
	require.NoError(t, err)
	fi, err := os.Stat(filepath.FromSlash(colimaPath))
	require.NotEmpty(t, fi)
	require.NoError(t, err)

	t.Logf("colimaPath=%s", colimaPath)
}

func TestGetUnixSocketAddreess(t *testing.T) {
	socketPath := getIncusSocketPath()
	raddr, err := net.ResolveUnixAddr("unix", socketPath)
	require.NoError(t, err)
	t.Logf("raddr=%#v", raddr)
}

func TestIsExistConfigFile(t *testing.T) {
	isExist, err := isExistConfigFile()
	require.NoError(t, err)
	require.True(t, isExist)
}
func TestIncusIsAvailable(t *testing.T) {
	is, err := isIncusAvailable()
	require.NoError(t, err)
	require.True(t, is)
}

func TestIncusIsDyltContainerExist (t *testing.T) {
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
