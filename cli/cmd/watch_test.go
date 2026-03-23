package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWatch(t *testing.T) {
	cmdName := "watch"
	cmdFlags := []string{}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, NewWatchCommand, cmdName, cmdFlags, cmdArgs, cmdString)
	require.IsType(t, &WatchCommand{}, cmd)
}

func TestWatchScript (t *testing.T) {
	// config get foo
	cmdName := "watch"
	subCmdName := "script"
	scriptKey := "scriptKey"
	targetPath := "targetPath"
	subCmdFlags := []string{}
	subCmdArgs := []string{scriptKey, targetPath}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := NewWatchCommand(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	subCmd := _TestSubcommandCreation[*WatchScriptCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, scriptKey, subCmd.ScriptKey)
	require.Equal(t, targetPath, subCmd.TargetPath)
}
func TestWatchSvc (t *testing.T) {
	// config get foo
	cmdName := "watch"
	subCmdName := "svc"
	name := "name"
	subCmdFlags := []string{}
	subCmdArgs := []string{name}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := NewWatchCommand(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	subCmd := _TestSubcommandCreation[*WatchSvcCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, name, subCmd.Name)
}
