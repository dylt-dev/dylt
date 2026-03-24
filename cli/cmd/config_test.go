package cmd

import (
	"fmt"
	"testing"

	"github.com/dylt-dev/dylt/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig (t *testing.T) {
	cmdName := "config"
	cmdFlags := []string{}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, NewConfigCommand, cmdName, cmdFlags, cmdArgs, cmdString)
	require.IsType(t, &ConfigCommand{}, cmd)
}

func TestConfigHelp(t *testing.T) {
	cmdName := "list"
	cmdFlags := []string{"--help"}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, NewConfigCommand, cmdName, cmdFlags, cmdArgs, cmdString)
	require.True(t, cmd.Help)
}

func TestConfigGet (t *testing.T) {
	// config get foo
	cmdName := "config"
	subCmdName := "get"
	subCmdFlags := []string{}
	subCmdArgs := []string{"foo"}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := NewConfigCommand(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	subCmd := _TestSubcommandCreation[*ConfigGetCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, "foo", subCmd.Key)
}


func TestConfigGetHelp(t *testing.T) {
	cmdName := "config"
	subCmdName := "get"
	subCmdFlags := []string{"--help"}
	subCmdArgs := []string{"foo"}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := NewConfigCommand(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	subCmd := _TestSubcommandCreation[*ConfigGetCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.True(t, subCmd.Help)
}

func TestConfigSet (t *testing.T) {
	// config set foo bar
	cmdName := "config"
	subCmdName := "set"
	subCmdFlags := []string{}
	subCmdArgs := Cmdline{"foo", "bar"}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := NewConfigCommand(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	subCmd := _TestSubcommandCreation[*ConfigSetCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, "foo", subCmd.Key)
	require.Equal(t, "bar", subCmd.Value)
}

func TestConfigSetHelp(t *testing.T) {
	cmdName := "config"
	subCmdName := "set"
	subCmdFlags := []string{"--help"}
	subCmdArgs := []string{"foo"}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := NewConfigCommand(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	subCmd := _TestSubcommandCreation[*ConfigSetCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.True(t, subCmd.Help)
}



func TestConfigShow (t *testing.T) {
	// config set foo bar
	cmdName := "config"
	subCmdName := "show"
	subCmdFlags := []string{}
	subCmdArgs := Cmdline{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := NewConfigCommand(cmdline, nil)
	cmd.HandleArgs()
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	subCmd := _TestSubcommandCreation[*ConfigShowCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.False(t, subCmd.Help)
}

func TestConfigShowHelp(t *testing.T) {
	cmdName := "config show"
	cmdFlags := []string{"--help"}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, NewConfigShowCommand, cmdName, cmdFlags, cmdArgs, cmdString)
	require.True(t, cmd.Help)
}

func TestRunConfigGet (t *testing.T) {
	key := "etcd-domain"
	err := RunConfigGet(key)
	require.Nil(t, err)
}

func TestConfigGetCmd (t *testing.T) {
	dyltPath := lib.GetAndValidateDyltPath(t)
	sCmdline := fmt.Sprintf("%s config get etcd-domain", dyltPath)
	lib.CheckRunCommandSuccess(sCmdline, t)
}


func TestRunConfigSet (t *testing.T) {
	key := "etcd-domain"
	val := "poo"
	err := RunConfigSet(key, val)
	assert.Nil(t, err)
}

func TestConfigSetCmd0 (t *testing.T) {
	sCmdline := "/Users/chris/src/dylt-dev/dylt/dylt config set etcd-domain MOO"
	CheckRunCommandSuccessNoOutput(sCmdline, t)
}

func TestConfigSetCmd1 (t *testing.T) {
	sCmdline := "/Users/chris/src/dylt-dev/dylt/dylt config set etcd-domain hello.dylt.dev"
	CheckRunCommandSuccessNoOutput(sCmdline, t)
}

// Created this test to figure out why I couldn't create Command objects with
// simple code. Turned out I was forgetting to Parse(). I was used to my 
// test_utils helper function doing it for me.
func TestCreateSubcommandRaw (t *testing.T) {
	var err error
	cmdline:= Cmdline{"config", "get", "--help"}
	cmd := NewConfigCommand(cmdline, nil)
	err = cmd.Parse()
	require.NoError(t, err)
	t.Logf("cmd=%v", cmd)
	t.Logf("cmd.CommandName()=%v", cmd.CommandName())
	require.NotNil(t, cmd)
	subCmd, err := cmd.CreateSubCommand()
	err = subCmd.Parse()
	require.NoError(t, err)
	require.NotNil(t, subCmd)
	require.True(t, subCmd.(*ConfigGetCommand).Help)
}


func TestRunConfigShow (t *testing.T) {
	err := RunConfigShow()
	assert.Nil(t, err)
}

func TestConfigShowCmd (t *testing.T) {
	sCmdline := "/Users/chris/src/dylt-dev/dylt/dylt config show"
	lib.CheckRunCommandSuccess(sCmdline, t)
}
