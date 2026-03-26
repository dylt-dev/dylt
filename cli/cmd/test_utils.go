package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dylt-dev/dylt/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// const PATH_Dylt = "/Users/chris/src/dylt-dev/dylt/dylt"

func CheckRunCommandSuccessNoOutput(sCmdlineArgs string, t *testing.T) error {
	dyltPath := lib.GetAndValidateDyltPath(t)
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
                                     fact func(Cmdline, Command) U,
                                     cmdName string,
                                     cmdFlags []string,
                                     cmdArgs []string,
                                     cmdString string) U {
	cmdline := NewCmdline(cmdName, cmdFlags, cmdArgs)
	t.Logf("cmdline=%v", cmdline)
	cmd := fact(cmdline, nil)
	err := cmd.HandleArgs()
	require.NoError(t, err)
	_TestCommandString(t, cmdString, cmd)
	require.Equal(t, cmdName, cmd.CommandName())
	args, is := cmd.Args()
	require.True(t, is)
	t.Logf("cmd.Args()=%v", args)
	require.Equal(t, Cmdline(cmdArgs), args)
	return cmd
}

func CreateCommandString(cmdName string, cmdArgs []string) string {
	return strings.Join(append([]string{cmdName}, cmdArgs...), " ")
}

func _TestCommandString(t *testing.T, targetCmdString string, cmd Command) {
	cmdString, flag := cmd.CommandString()
	require.True(t, flag)
	require.Equal(t, targetCmdString, cmdString)
}

// _TestSubCommandAndArgs
//
// This function does *not* test subcommand creation. It just confirms
// expectations for the Command values that are used for subcommand creation
//
// Params
//    cmd               parent command, loaded with a subCmd-producing Cmdline
//	  targetSubCommand	the expected subcommand name
//    targetSubArgs		the expected subcommand args
func _TestSubCommandAndArgs(t *testing.T,
	                        cmd Command,
							targetSubCommand string,
							targetSubCmdline []string) {
	subCommand, flag := cmd.SubCommand()
	require.Equal(t, targetSubCommand, subCommand)
	require.True(t, flag)
	subCmdline, flag := cmd.Args()
	require.Equal(t, Cmdline(targetSubCmdline), subCmdline)
	require.True(t, flag)
}

func _TestParentCommand(t *testing.T,
                        cmd Command,
                        cmdName string,
                        cmdArgs Cmdline) {
	err := cmd.HandleArgs()
	require.NoError(t, err)
	_TestCommandString(t, cmdName, cmd)
	args, is := cmd.Args()
	require.True(t, is)
	require.Equal(t, cmdArgs, args)
}

// _TestSubcommandCreation
//
// Params
//
//    cmd               parent command, loaded with a subCmd-producing Cmdline
//    subName			name of the subcommand
//    subCmdFlags		flags to pass to the subcommand
//    subArgs			arguments to pass to the subcommand
//    targetCmdString	expected command string 
func _TestSubcommandCreation[TCmd Command](t *testing.T,
                                           cmd Command,
                                           subName string,
                                           subCmdFlags []string,
                                           targetSubArgs []string,
                                           targetCmdString string) TCmd {
	var err error
	subCmdline := NewCmdline(subName, subCmdFlags, targetSubArgs)
	_TestSubCommandAndArgs(t, cmd, subName, subCmdline)
	// Create subcommand, HandleArg(), confirm return type
	subCmdRaw, err := cmd.CreateSubCommand()
	require.NoError(t, err)
	require.NotNil(t, subCmdRaw)
	subCmd, ok := subCmdRaw.(TCmd)
	require.True(t, ok)
	err = subCmd.HandleArgs()
	require.NoError(t, err)
	// name
	commandName := subCmd.CommandName()
	require.Equal(t, subName, commandName)
	// args
	subArgs, is := subCmd.Args()
	require.True(t, is)
	require.Equal(t, Cmdline(targetSubArgs), subArgs)
	// cmdString
	cmdString, is := subCmd.CommandString()
	require.True(t, is)
	require.Equal(t, targetCmdString, cmdString)

	return subCmd
}
