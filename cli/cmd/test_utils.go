package cmd

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/dylt-dev/dylt/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// const PATH_Dylt = "/Users/chris/src/dylt-dev/dylt/dylt"

func CheckRunCommandSuccess (sCmdlineArgs string, t *testing.T) error {
	dyltPath := GetAndValidateDyltPath(t)
	rc, stdout, err := lib.RunCommand(dyltPath, strings.Fields(sCmdlineArgs)...)
	require.Equal(t, 0, rc)
	require.NotEmpty(t, stdout)
	require.Nil(t, err)
	t.Log(string(stdout))
	return err
}


func CheckRunCommandSuccessNoOutput (sCmdlineArgs string, t *testing.T) error {
	dyltPath := GetAndValidateDyltPath(t)
	sCmdline := fmt.Sprintf("%s %s", dyltPath, sCmdlineArgs)
	var cmdline Cmdline = strings.Split(sCmdline, " ")
	rc, stdout, err := lib.RunCommand(cmdline.Command(), cmdline.Args()...)
	assert.Equal(t, 0, rc)
	assert.Empty(t, stdout)
	assert.Nil(t, err)
	t.Log(string(stdout))
	return err
}

func GetAndValidateDyltPath (t *testing.T) string {
	envName := "DYLT_PATH"
	dyltPath, is := os.LookupEnv(envName)
	if !is {
		t.Skipf("%s not set", envName)
	}
	_, err := os.Stat(dyltPath)
	if !os.IsNotExist(err) {
		t.Skipf("dylt path not found: %s", dyltPath)
	}
	return dyltPath
}

func _TestSubCommandAndArgs (t *testing.T, cmd Command, targetSubCommand string, targetSubArgs Cmdline) {
	subCommand, flag := cmd.SubCommand()
	require.Equal(t, targetSubCommand, subCommand)
	require.True(t, flag)
	subArgs, flag := cmd.SubArgs()
	require.Equal(t, targetSubArgs, subArgs)	
	require.True(t, flag)
}

func _TestSubcommandCreation[TCmd Command] (t *testing.T,
	                          cmd SuperCommand,
							  subName string,
							  subArgs Cmdline,
							  targetCmdString string) Command {
	var err error
	_TestSubCommandAndArgs(t, cmd, subName, subArgs)
	// Create + parse subcommand
	subCmd, err := cmd.CreateSubCommand()
	require.NoError(t, err)
	err = subCmd.HandleArgs()
	require.NoError(t, err)
	// subCommand type
	// require.IsType(t, &TCmd{}, subCmd)
	_, ok := subCmd.(TCmd)
	require.True(t, ok)
	// cmdString
	cmdString, flag := subCmd.GetCommandString()
	require.True(t, flag)
	require.Equal(t, targetCmdString, cmdString)
	// name
	commandName := subCmd.CommandName()
	require.Equal(t, subName, commandName)
	// args
	args, flag := subCmd.Args()
	require.True(t, flag)
	require.Equal(t, subArgs, args)

	return subCmd
}