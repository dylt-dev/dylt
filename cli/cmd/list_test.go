package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunList (t *testing.T) {
	err := RunList()
	require.Nil(t, err)
}

func TestListCmd0 (t *testing.T) {
	if os.Getenv("DYLT_TEST_SYSTEST") == "" {
		t.Skip()
	}
	dyltPath, ok := os.LookupEnv("DYLT_EXE_PATH")
	require.True(t, ok)
	require.NotEmpty(t, dyltPath)
	cmd := fmt.Sprintf("%s list", dyltPath)
	err := CheckRunCommandSuccess(cmd, t)
	require.Nil(t, err)
}

func TestListHandleArgs_None (t *testing.T) {
	args := []string{}
	cmd := NewListCommand(args)
	err := cmd.HandleArgs()
	require.NoError(t, err)
	require.Empty(t, cmd.SubArgs)
	require.Empty(t, cmd.SubCommand)
}