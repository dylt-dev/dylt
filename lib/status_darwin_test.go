package lib

import (
	"os/exec"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

func TestColimaGetSocketPath (t *testing.T) {
	path := getColimaSocketPath()
	require.NotEmpty(t, path)
	t.Logf("path=%s", path)
}

func TestColimaGetStatus (t *testing.T) {
	common.InitLogging()
	
	cmd := exec.Command("colima", "status")
	bufferStdout, bufferStderr, err := runWithStdoutAndStderr(cmd)
	
	t.Log(bufferStdout.String())
	t.Log(bufferStderr.String())
	
	require.NoError(t, err)
	require.NotEmpty(t, bufferStdout)
	require.NotEmpty(t, bufferStderr)
}

func TestColimaIsActive (t *testing.T) {
	flag, err := isColimaActive()
	require.NoError(t, err)
	require.True(t, flag)
}
func TestColimaStartIncus (t *testing.T) {
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
