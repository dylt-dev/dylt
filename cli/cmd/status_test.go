package cmd

import (
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

func TestStatus (t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	cmdName := "status"
	cmdFlags := []string{}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, StatusCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString).(*BaseCommand[StatusOpts])
	require.IsType(t, &BaseCommand[StatusOpts]{}, cmd)
	require.False(t, cmd.Help())
}

func TestStatusHelp (t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	cmdName := "status"
	cmdFlags := []string{"--help"}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, StatusCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString).(*BaseCommand[StatusOpts])
	require.IsType(t, &BaseCommand[StatusOpts]{}, cmd)
	require.True(t, cmd.Help())
}