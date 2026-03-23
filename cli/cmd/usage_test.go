package cmd

import (
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

func setup (t *testing.T) func (t *testing.T) {
	t.Log("setup() ...")
	common.InitLogging()
	return teardown
}

func teardown (t *testing.T) {
	t.Log("teardown() ...")
}

func TestMainUsage (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)

	cmdline := NewCmdline("dylt", []string{}, []string{})
	cmd := NewMainCommand(cmdline, nil)
	require.NotNil(t, cmd)
	err := cmd.Run()
	require.NoError(t, err)
}

func TestConfigUsage (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdline := Cmdline{"dylt", "config"}

	cmd := NewMainCommand(cmdline, nil)
	cmd.HandleArgs()
	subCmd, err := cmd.CreateSubCommand()
	require.NoError(t, err)
	require.NotNil(t, subCmd)
	err = subCmd.Run()
	require.NoError(t, err)
}

func TestHostUsage (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdline := Cmdline{"dylt", "host"}

	cmd := NewMainCommand(cmdline, nil)
	cmd.HandleArgs()
	subCmd, err := cmd.CreateSubCommand()
	require.NoError(t, err)
	require.NotNil(t, subCmd)
	err = subCmd.Run()
	require.NoError(t, err)
}

func TestMiscUsage (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdline := Cmdline{"dylt", "misc"}

	cmd := NewMainCommand(cmdline, nil)
	cmd.HandleArgs()
	subCmd, err := cmd.CreateSubCommand()
	require.NoError(t, err)
	require.NotNil(t, subCmd)
	err = subCmd.Run()
	require.NoError(t, err)
}

func TestVmUsage (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdline := Cmdline{"dylt", "vm"}

	cmd := NewMainCommand(cmdline, nil)
	cmd.HandleArgs()
	subCmd, err := cmd.CreateSubCommand()
	require.NoError(t, err)
	require.NotNil(t, subCmd)
	err = subCmd.Run()
	require.NoError(t, err)
}

func TestWatchUsage (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdline := Cmdline{"dylt", "watch"}

	cmd := NewMainCommand(cmdline, nil)
	cmd.HandleArgs()
	subCmd, err := cmd.CreateSubCommand()
	require.NoError(t, err)
	require.NotNil(t, subCmd)
	err = subCmd.Run()
	require.NoError(t, err)
}