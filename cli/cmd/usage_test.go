package cmd

import (
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

func TestConfigUsage (t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	cmdline := Cmdline{"dylt", "config"}

	cmd := MainCommandF.New(cmdline, nil)
	cmd.HandleArgs()
	subCmd, err := cmd.CreateSubCommand()
	require.NoError(t, err)
	require.NotNil(t, subCmd)
	err = subCmd.Run()
	require.NoError(t, err)
}

func TestHostUsage (t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	cmdline := Cmdline{"dylt", "host"}

	cmd := MainCommandF.New(cmdline, nil)
	cmd.HandleArgs()
	subCmd, err := cmd.CreateSubCommand()
	require.NoError(t, err)
	require.NotNil(t, subCmd)
	err = subCmd.Run()
	require.NoError(t, err)
}

func TestMiscUsage (t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	cmdline := Cmdline{"dylt", "misc"}

	cmd := MainCommandF.New(cmdline, nil)
	cmd.HandleArgs()
	subCmd, err := cmd.CreateSubCommand()
	require.NoError(t, err)
	require.NotNil(t, subCmd)
	err = subCmd.Run()
	require.NoError(t, err)
}

func TestVmUsage (t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	cmdline := Cmdline{"dylt", "vm"}

	cmd := MainCommandF.New(cmdline, nil)
	cmd.HandleArgs()
	subCmd, err := cmd.CreateSubCommand()
	require.NoError(t, err)
	require.NotNil(t, subCmd)
	err = subCmd.Run()
	require.NoError(t, err)
}

func TestWatchUsage (t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	cmdline := Cmdline{"dylt", "watch"}

	cmd := MainCommandF.New(cmdline, nil)
	cmd.HandleArgs()
	subCmd, err := cmd.CreateSubCommand()
	require.NoError(t, err)
	require.NotNil(t, subCmd)
	err = subCmd.Run()
	require.NoError(t, err)
}