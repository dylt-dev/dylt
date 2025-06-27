package lib

import (
	"os"
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

func TestIsExistConfigFile (t *testing.T) {
	isExist, err := isExistConfigFile()
	require.NoError(t, err)
	require.True(t, isExist)
}
