package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWatch(t *testing.T) {
	cmdName := "watch"
	cmdFlags := []string{}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, WatchCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString)
	require.IsType(t, &WatchCommand{}, cmd)
}


func TestWatchHelp (t *testing.T) {
	cmdName := "watch"
	cmdFlags := []string{"--help"}
	cmdArgs := []string{}
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t, WatchCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString).(*WatchCommand)
	require.IsType(t, &WatchCommand{}, cmd)
	require.True(t, cmd.Help)
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
	cmd := WatchCommandF.New(cmdline, nil)
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


func TestWatchScriptHelp(t *testing.T) {
	cmdName := "watch"
	subCmdName := "script"
	subCmdFlags := []string{"--help"}
	subCmdArgs := []string{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := WatchCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test  subcommand
	subCmd := _TestSubcommandCreation[*WatchScriptCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.True(t, subCmd.Help)
}

func TestWatchSvc (t *testing.T) {
	// config get foo
	cmdName := "watch"
	subCmdName := "svc"
	name := "name"
	subCmdFlags := []string{}
	subCmdArgs := []string{name}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := WatchCommandF.New(cmdline, nil)
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


func TestWatchSvcHelp(t *testing.T) {
	cmdName := "watch"
	subCmdName := "svc"
	subCmdFlags := []string{"--help"}
	subCmdArgs := []string{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := WatchCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test  subcommand
	subCmd := _TestSubcommandCreation[*WatchSvcCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.True(t, subCmd.Help)
}