package lib

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func RunAndGetOutput(cmdLine string) (string, error) {
	cmdAndArgs := strings.Fields(cmdLine)
	cmdName := cmdAndArgs[0]
	args := cmdAndArgs[1:]
	cmd := exec.Command(cmdName, args...)
	out, err := cmd.Output()
	s := string(out)
	return s, err
}

func RunAndTestCommand(t *testing.T, cmdline string) {
	out, err := RunAndGetOutput(cmdline)
	require.NoError(t, err)
	t.Log(string(out))
}
