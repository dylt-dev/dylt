package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateCommandParams (cmdName string, subCmdName string, subCmdArgs Cmdline) (cmdline Cmdline, cmdArgs Cmdline, subCmdString string) {
	cmdArgs = append(Cmdline{subCmdName}, subCmdArgs...)
	cmdline = append(Cmdline{cmdName}, cmdArgs...)
	subCmdString = strings.Join(cmdline[0:2], " ")
	return
}

func TestMain (t *testing.T) {
	var err error
	cmdline := Cmdline{"dylt"}

	cmd := NewMainCommand(cmdline)
	err = cmd.Parse()
	require.NoError(t, err)
	_TestCommandString(t, cmd, "dylt")
	require.False(t, cmd.Help)
}


func TestMainHelp (t *testing.T) {
	var err error
	cmdline := Cmdline{"dylt", "--help"}

	cmd := NewMainCommand(cmdline)
	err = cmd.HandleArgs()
	require.NoError(t, err)
	require.True(t, cmd.Help)
	cmdString, _ := cmd.GetCommandString()
	require.Equal(t, "dylt", cmdString)
}


// Test that `dylt call foo` produces a `call` subcommand
// where name == "call" and args == ["foo"]
func TestMainSubCall (t *testing.T) {
	// Create + parse main command with "dylt call foo"
	cmdline := Cmdline{"dylt", "call", "foo"}
	cmd := NewMainCommand(cmdline)
	err := cmd.Parse()
	require.NoError(t, err)
	_TestSubcommandCreation[*CallCommand](t,
		cmd,
		"call",
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

func TestMainSubConfig (t *testing.T) {
	// Create + parse main command with "dylt config get foo"
	cmdline := Cmdline{"dylt", "config", "get", "foo"}
	cmd := NewMainCommand(cmdline)
	err := cmd.Parse()
	// subcommand == "config", subArgs == {"get", "foo"}
	_TestSubCommandAndArgs(t, cmd, "config", []string{"get", "foo"})
	require.NoError(t, err)
	_TestSubcommandCreation[*ConfigCommand](t,
		cmd,
		"config",
		Cmdline{"get", "foo"},
		"dylt config",
	)
}

func TestMainSubGet (t *testing.T) {
	// Create + parse main command with "dylt get foo"
	cmdline := Cmdline{"dylt", "get", "foo"}
	cmd := NewMainCommand(cmdline)
	err := cmd.Parse()
	// subcommand == "get", subArgs == {"foo"}
	_TestSubCommandAndArgs(t, cmd, "get", []string{"foo"})
	require.NoError(t, err)
	_TestSubcommandCreation[*GetCommand](t,
		cmd,
		"get",
		Cmdline{"foo"},
		"dylt get",
	)
}

func TestMainSubHost (t *testing.T) {
	// Create + parse main command with "dylt get foo"
	cmdline := Cmdline{"dylt", "host", "init", "2000", "2000"}
	cmd := NewMainCommand(cmdline)
	err := cmd.Parse()
	// subcommand == "host", subArgs == {"foo"}
	_TestSubCommandAndArgs(t, cmd, "host", []string{"init", "2000", "2000"})
	require.NoError(t, err)
	_TestSubcommandCreation[*HostCommand](t,
		cmd,
		"host",
		Cmdline{"init", "2000", "2000"},
		"dylt host",
	)
}

func TestMainSubInit (t *testing.T) {
	// Create + parse main command with "dylt get foo"
	subCmdName := "init"
	subCmdArgs := []string{"etcdDomain"}
	cmdline := append([]string{"dylt", subCmdName}, subCmdArgs...)
	cmd := NewMainCommand(cmdline)
	err := cmd.Parse()
	// subcommand == "host", subArgs == {"foo"}
	_TestSubCommandAndArgs(t, cmd, subCmdName, subCmdArgs)
	require.NoError(t, err)
	_TestSubcommandCreation[*InitCommand](t,
		cmd,
		subCmdName,
		subCmdArgs,
		"dylt init",
	)
}

func TestMainSubList (t *testing.T) {
	// Create + parse main command with "dylt list etcdDomain"
	subCmdName := "list"
	subCmdArgs := []string{}
	cmdline := append([]string{"dylt", subCmdName}, subCmdArgs...)
	cmdString := strings.Join(cmdline[0:2], " ")
	cmd := NewMainCommand(cmdline)
	err := cmd.Parse()
	// subcommand == "host", subArgs == {"foo"}
	_TestSubCommandAndArgs(t, cmd, subCmdName, subCmdArgs)
	require.NoError(t, err)
	_TestSubcommandCreation[*ListCommand](t,
		cmd,
		subCmdName,
		subCmdArgs,
		cmdString,
	)
}

func TestMainSubStatus (t *testing.T) {
	// Create + parse main command with "dylt misc lookup hostname"
	subCmdName := "status"
	subCmdArgs := []string{}
	cmdline := append([]string{"dylt", subCmdName}, subCmdArgs...)
	cmdString := strings.Join(cmdline[0:2], " ")
	cmd := NewMainCommand(cmdline)
	err := cmd.Parse()
	// subcommand & subArgs
	_TestSubCommandAndArgs(t, cmd, subCmdName, subCmdArgs)
	require.NoError(t, err)
	_TestSubcommandCreation[*StatusCommand](t,
		cmd,
		subCmdName,
		subCmdArgs,
		cmdString,
	)
}

func TestMainSubMisc (t *testing.T) {
	// Create + parse main command with "dylt misc lookup hostname"
	subCmdName := "misc"
	subCmdArgs := []string{"lookup", "hostname"}
	cmdline := append([]string{"dylt", subCmdName}, subCmdArgs...)
	cmdString := strings.Join(cmdline[0:2], " ")
	cmd := NewMainCommand(cmdline)
	err := cmd.Parse()
	// subcommand & subArgs
	_TestSubCommandAndArgs(t, cmd, subCmdName, subCmdArgs)
	require.NoError(t, err)
	_TestSubcommandCreation[*MiscCommand](t,
		cmd,
		subCmdName,
		subCmdArgs,
		cmdString,
	)
}

func TestMainSubVm (t *testing.T) {
	// Create + parse main command with "dylt vm get name"
	subCmdName := "vm"
	subCmdArgs := []string{"get", "name"}
	cmdline := append([]string{"dylt", subCmdName}, subCmdArgs...)
	cmdString := strings.Join(cmdline[0:2], " ")
	cmd := NewMainCommand(cmdline)
	err := cmd.Parse()
	// subcommand & subArgs
	_TestSubCommandAndArgs(t, cmd, subCmdName, subCmdArgs)
	require.NoError(t, err)
	_TestSubcommandCreation[*VmCommand](t,
		cmd,
		subCmdName,
		subCmdArgs,
		cmdString,
	)
}

func TestMainSubWatch (t *testing.T) {
	// dylt watch script scriptPath
	cmdName := "dylt"
	subCmdName := "watch"
	subCmdArgs := Cmdline{"script", "scriptPath"}
	cmdline := append([]string{cmdName, subCmdName}, subCmdArgs...)
	subCmdString := strings.Join(cmdline[0:2], " ")
	cmd := NewMainCommand(cmdline)
	err := cmd.Parse()
	require.NoError(t, err)
	// subcommand & subArgs
	_TestSubCommandAndArgs(t, cmd, subCmdName, subCmdArgs)
	require.NoError(t, err)
	_TestSubcommandCreation[*WatchCommand](t,
		cmd,
		subCmdName,
		subCmdArgs,
		subCmdString,
	)
}

func _TestCommandString (t *testing.T, cmd Command, targetCmdString string) {
	cmdString, flag := cmd.GetCommandString()
	require.True(t, flag)
	require.Equal(t, targetCmdString, cmdString)
}

func _TestParentCommand (t *testing.T,
	                     cmd SuperCommand,
						 cmdName string,
						 cmdArgs Cmdline) {
	err := cmd.Parse()
	require.NoError(t, err)
	_TestCommandString(t, cmd, cmdName)
	args, is := cmd.Args()
	require.True(t, is)
	require.Equal(t, cmdArgs, args)
}