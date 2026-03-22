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

func CheckRunCommandSuccess(sCmdlineArgs string, t *testing.T) error {
	dyltPath := GetAndValidateDyltPath(t)
	rc, stdout, err := lib.RunCommand(dyltPath, strings.Fields(sCmdlineArgs)...)
	require.Equal(t, 0, rc)
	require.NotEmpty(t, stdout)
	require.Nil(t, err)
	t.Log(string(stdout))
	return err
}

func CheckRunCommandSuccessNoOutput(sCmdlineArgs string, t *testing.T) error {
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

func CreateAndTestCommand[U Command](t *testing.T,
                                     fact func(Cmdline, SuperCommand) U,
                                     cmdName string,
                                     cmdFlags []string,
                                     cmdArgs []string,
                                     cmdString string) U {
	cmdline := NewCmdline(cmdName, cmdFlags, cmdArgs)
	cmd := fact(cmdline, nil)
	err := cmd.HandleArgs()
	require.NoError(t, err)
	_TestCommandString(t, cmdString, cmd)
	require.Equal(t, cmdName, cmd.CommandName())
	args, is := cmd.Args()
	require.True(t, is)
	require.Equal(t, Cmdline(cmdArgs), args)
	return cmd
}

func CreateCommandString(cmdName string, cmdArgs []string) string {
	return strings.Join(append([]string{cmdName}, cmdArgs...), " ")
}

func GetAndValidateDyltPath(t *testing.T) string {
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

func _TestCommandString(t *testing.T, targetCmdString string, cmd Command) {
	cmdString, flag := cmd.CommandString()
	require.True(t, flag)
	require.Equal(t, targetCmdString, cmdString)
}

func _TestSubCommandAndArgs(t *testing.T,
	                        cmd Command,
							targetSubCommand string,
							targetSubArgs Cmdline) {
	subCommand, flag := cmd.SubCommand()
	require.Equal(t, targetSubCommand, subCommand)
	require.True(t, flag)
	args, flag := cmd.Args()
	require.Equal(t, targetSubArgs, args)
	require.True(t, flag)
}

func _TestParentCommand(t *testing.T,
                        cmd SuperCommand,
                        cmdName string,
                        cmdArgs Cmdline) {
	err := cmd.HandleArgs()
	require.NoError(t, err)
	_TestCommandString(t, cmdName, cmd)
	args, is := cmd.Args()
	require.True(t, is)
	require.Equal(t, cmdArgs, args)
}

func _TestSubcommandCreation[TCmd Command](t *testing.T,
                                           cmd SuperCommand,
                                           subName string,
                                           subCmdFlags []string,
                                           subArgs Cmdline,
                                           targetCmdString string) Command {
	var err error
	targetSubArgs := NewCmdline(subName, subCmdFlags, subArgs)
	_TestSubCommandAndArgs(t, cmd, subName, targetSubArgs)
	// Create + parse subcommand
	subCmd, err := cmd.CreateSubCommand()
	require.NoError(t, err)
	require.NotNil(t, subCmd)
	t.Logf("subCmd.(*BaseCommand).Cmdline=%v", subCmd.(TCmd).CommandLine())
	err = subCmd.HandleArgs()
	require.NoError(t, err)
	// subCommand type
	// require.IsType(t, &TCmd{}, subCmd)
	_, ok := subCmd.(TCmd)
	require.True(t, ok)
	// cmdString
	cmdString, is := subCmd.CommandString()
	require.True(t, is)
	require.Equal(t, targetCmdString, cmdString)
	// name
	commandName := subCmd.CommandName()
	require.Equal(t, subName, commandName)
	// args
	args, is := subCmd.Args()
	require.True(t, is)
	require.Equal(t, subArgs, args)

	return subCmd
}
