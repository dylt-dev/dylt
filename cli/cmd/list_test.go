package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/dylt-dev/dylt/api"
	"github.com/dylt-dev/dylt/lib"
	"github.com/stretchr/testify/require"
)

func TestList (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdName := "list"
	cmdFlags := []string{}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, ListCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString).(*BaseCommand[ListOpts])
	require.NotNil(t, cmd)
}

func TestListHelp(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdName := "list"
	cmdFlags := []string{"--help"}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, ListCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString).(*BaseCommand[ListOpts])
	require.NotNil(t, cmd)
	require.True(t, cmd.Help())
}

func TestRunList (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	err := api.RunList()
	require.Nil(t, err)
}

func TestListCmd0 (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	if os.Getenv("DYLT_TEST_SYSTEST") == "" {
		t.Skip()
	}
	dyltPath, ok := os.LookupEnv("DYLT_EXE_PATH")
	require.True(t, ok)
	require.NotEmpty(t, dyltPath)
	cmd := fmt.Sprintf("%s list", dyltPath)
	err := lib.CheckRunCommandSuccess(cmd, t)
	require.Nil(t, err)
}

func TestListHandleArgs_None (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	args := []string{}
	cmd := ListCommandF.New(args, nil)
	err := cmd.HandleArgs()
	require.NoError(t, err)
	subArgs, _ := cmd.SubArgs()
	require.Empty(t, subArgs)
	subCommand, _ := cmd.SubCommand()
	require.Empty(t, subCommand)
}