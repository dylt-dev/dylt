package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStatus (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdName := "status"
	cmdFlags := []string{}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, StatusCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString).(*StatusCommand)
	require.IsType(t, &StatusCommand{}, cmd)
	require.False(t, cmd.Help)
}

func TestStatusHelp (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdName := "status"
	cmdFlags := []string{"--help"}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, StatusCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString).(*StatusCommand)
	require.IsType(t, &StatusCommand{}, cmd)
	require.True(t, cmd.Help)
}