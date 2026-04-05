package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWatch(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdName := "watch"
	cmdFlags := []string{}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, WatchCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString)
	require.IsType(t, &BaseCommand[WatchOpts]{}, cmd)
}


func TestWatchHelp (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdName := "watch"
	cmdFlags := []string{"--help"}
	cmdArgs := []string{}
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t, WatchCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString).(*BaseCommand[WatchOpts])
	require.IsType(t, &BaseCommand[WatchOpts]{}, cmd)
	require.True(t, cmd.Help())
}

func TestWatchScript (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
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
	subCmd := _TestSubcommandCreation[*BaseCommand[WatchScriptOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, scriptKey, subCmd.opts.ScriptKey)
	require.Equal(t, targetPath, subCmd.opts.TargetPath)
}


func TestWatchScriptHelp(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdName := "watch"
	subCmdName := "script"
	subCmdFlags := []string{"--help"}
	subCmdArgs := []string{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := WatchCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test  subcommand
	subCmd := _TestSubcommandCreation[*BaseCommand[WatchScriptOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.True(t, subCmd.Help())
}

func TestWatchSvc (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
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
	subCmd := _TestSubcommandCreation[*BaseCommand[WatchSvcOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, name, subCmd.opts.Name)
}


func TestWatchSvcHelp(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdName := "watch"
	subCmdName := "svc"
	subCmdFlags := []string{"--help"}
	subCmdArgs := []string{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := WatchCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test  subcommand
	subCmd := _TestSubcommandCreation[*BaseCommand[WatchSvcOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.True(t, subCmd.Help())
}