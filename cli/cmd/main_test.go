package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Returns
//     cmdline	    The command line for the command+subcommand invocation
//     cmdArgs      The arguments for the command
//     subCmdString The command string for the full command: cmd, subCmd, subCmdArgs
func CreateCommandParams(cmdName string,
	                     subCmdName string,
						 subCmdFlags []string,
						 subCmdArgs []string) (cmdline Cmdline, cmdArgs Cmdline, cmdString string) {
	subCmdline := NewCmdline(subCmdName, subCmdFlags, subCmdArgs)
	cmdline = NewCmdline(cmdName, []string{}, subCmdline)
	cmdArgs = cmdline.Args()
	cmdString = strings.Join(cmdline[0:2], " ")
	return
}

func TestMain(t *testing.T) {
	cmdName := "dylt"
	cmdFlags := []string{}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, NewMainCommand, cmdName, cmdFlags, cmdArgs, cmdString)
	require.False(t, cmd.Help)
}

func TestMainHelp(t *testing.T) {
	cmdName := "dylt"
	cmdFlags := []string{"--help"}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, NewMainCommand, cmdName, cmdFlags, cmdArgs, cmdString)
	require.True(t, cmd.Help)
}

// Test that `dylt call foo` produces a `call` subcommand
// where name == "call" and args == ["foo"]
func TestMainSubCall(t *testing.T) {
	// Create + parse main command with "dylt call foo"
	cmdline := Cmdline{"dylt", "call", "foo"}
	cmd := NewMainCommand(cmdline, nil)
	err := cmd.Parse()
	require.NoError(t, err)
	_TestSubcommandCreation[*CallCommand](t,
		cmd,
		"call",
		[]string{},
		Cmdline{"foo"},
		"dylt call",
	)
	// _TestSubCommandAndArgs(t, cmd, "call", []string{"foo"})
	// // Create + parse subcommand
	// subCmd, err := cmd.CreateSubCommand()
	// require.NoError(t, err)
	// err = subCmd.Parse()
	// require.NoError(t, err)
	// // subCommand type = CallCommand
	// require.IsType(t, &CallCommand{}, subCmd)
	// // cmdString == "dylt call"
	// cmdString, flag := subCmd.GetCommandString()
	// require.True(t, flag)
	// require.Equal(t, "dylt call", cmdString)
	// // name == "call"
	// commandName := subCmd.CommandName()
	// require.Equal(t, "call", commandName)
	// // args == ["foo"]
	// args, flag := subCmd.Args()
	// require.True(t, flag)
	// require.Equal(t, 1, len(args))
	// require.Equal(t, "foo", args[0])
}

func TestMainSubConfig(t *testing.T) {
	// Create + parse main command with "dylt config get foo"
	cmdline := Cmdline{"dylt", "config", "get", "foo"}
	cmd := NewMainCommand(cmdline, nil)
	err := cmd.Parse()
	// subcommand == "config", subArgs == {"get", "foo"}
	_TestSubCommandAndArgs(t, cmd, "config", []string{"get", "foo"})
	require.NoError(t, err)
	_TestSubcommandCreation[*ConfigCommand](t,
		cmd,
		"config",
		[]string{},
		Cmdline{"get", "foo"},
		"dylt config",
	)
}

func TestMainSubGet(t *testing.T) {
	// Create + parse main command with "dylt get foo"
	cmdline := Cmdline{"dylt", "get", "foo"}
	cmd := NewMainCommand(cmdline, nil)
	err := cmd.Parse()
	// subcommand == "get", subArgs == {"foo"}
	_TestSubCommandAndArgs(t, cmd, "get", []string{"foo"})
	require.NoError(t, err)
	_TestSubcommandCreation[*GetCommand](t,
		cmd,
		"get",
		[]string{},
		Cmdline{"foo"},
		"dylt get",
	)
}

func TestMainSubHost(t *testing.T) {
	// Create + parse main command with "dylt get foo"
	cmdline := Cmdline{"dylt", "host", "init", "--gid", "1000", "--uid", "2000"}
	cmd := NewMainCommand(cmdline, nil)
	err := cmd.Parse()
	// subcommand == "host", subArgs == {"foo"}
	_TestSubCommandAndArgs(t, cmd, "host", []string{"init", "--gid", "1000", "--uid", "2000"})
	require.NoError(t, err)
	_TestSubcommandCreation[*HostCommand](t,
		cmd,
		"host",
		[]string{},
		Cmdline{"init", "--gid", "1000", "--uid", "2000"},
		"dylt host",
	)
}

func TestMainSubInit(t *testing.T) {
	// Create + parse main command with "dylt get foo"
	subCmdName := "init"
	subCmdFlags := []string{}
	subCmdArgs := []string{"etcdDomain"}
	cmdline := append([]string{"dylt", subCmdName}, subCmdArgs...)
	cmd := NewMainCommand(cmdline, nil)
	err := cmd.Parse()
	// subcommand == "host", subArgs == {"foo"}
	_TestSubCommandAndArgs(t, cmd, subCmdName, subCmdArgs)
	require.NoError(t, err)
	_TestSubcommandCreation[*InitCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		"dylt init",
	)
}

func TestMainSubList(t *testing.T) {
	// Create + parse main command with "dylt list etcdDomain"
	subCmdName := "list"
	subCmdFlags := []string{}
	subCmdArgs := []string{}
	cmdline := append([]string{"dylt", subCmdName}, subCmdArgs...)
	cmdString := strings.Join(cmdline[0:2], " ")
	cmd := NewMainCommand(cmdline, nil)
	err := cmd.Parse()
	// subcommand == "host", subArgs == {"foo"}
	_TestSubCommandAndArgs(t, cmd, subCmdName, subCmdArgs)
	require.NoError(t, err)
	_TestSubcommandCreation[*ListCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		cmdString,
	)
}

func TestMainSubStatus(t *testing.T) {
	// Create + parse main command with "dylt misc lookup hostname"
	subCmdName := "status"
	subCmdFlags := []string{}
	subCmdArgs := []string{}
	cmdline := append([]string{"dylt", subCmdName}, subCmdArgs...)
	cmdString := strings.Join(cmdline[0:2], " ")
	cmd := NewMainCommand(cmdline, nil)
	err := cmd.Parse()
	// subcommand & subArgs
	_TestSubCommandAndArgs(t, cmd, subCmdName, subCmdArgs)
	require.NoError(t, err)
	_TestSubcommandCreation[*StatusCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		cmdString,
	)
}

func TestMainSubMisc(t *testing.T) {
	// Create + parse main command with "dylt misc lookup hostname"
	subCmdName := "misc"
	subCmdFlags := []string{}
	subCmdArgs := []string{"lookup", "hostname"}
	cmdline := append([]string{"dylt", subCmdName}, subCmdArgs...)
	cmdString := strings.Join(cmdline[0:2], " ")
	cmd := NewMainCommand(cmdline, nil)
	err := cmd.Parse()
	// subcommand & subArgs
	_TestSubCommandAndArgs(t, cmd, subCmdName, subCmdArgs)
	require.NoError(t, err)
	_TestSubcommandCreation[*MiscCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		cmdString,
	)
}

func TestMainSubVm(t *testing.T) {
	// Create + parse main command with "dylt vm get name"
	subCmdName := "vm"
	subCmdFlags := []string{}
	subCmdArgs := []string{"get", "name"}
	cmdline := append([]string{"dylt", subCmdName}, subCmdArgs...)
	cmdString := strings.Join(cmdline[0:2], " ")
	cmd := NewMainCommand(cmdline, nil)
	err := cmd.Parse()
	// subcommand & subArgs
	_TestSubCommandAndArgs(t, cmd, subCmdName, subCmdArgs)
	require.NoError(t, err)
	_TestSubcommandCreation[*VmCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		cmdString,
	)
}

func TestMainSubWatch(t *testing.T) {
	// dylt watch script scriptPath
	cmdName := "dylt"
	subCmdName := "watch"
	subCmdFlags := []string{}
	subCmdArgs := Cmdline{"script", "scriptPath"}
	cmdline := append([]string{cmdName, subCmdName}, subCmdArgs...)
	subCmdString := strings.Join(cmdline[0:2], " ")
	cmd := NewMainCommand(cmdline, nil)
	err := cmd.Parse()
	require.NoError(t, err)
	// subcommand & subArgs
	_TestSubCommandAndArgs(t, cmd, subCmdName, subCmdArgs)
	require.NoError(t, err)
	_TestSubcommandCreation[*WatchCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
}
