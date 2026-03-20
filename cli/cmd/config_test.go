package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig (t *testing.T) {
	var err error
	cmdline := Cmdline{"config"}

	cmd := NewConfigCommand(cmdline, nil)
	err = cmd.Parse()
	require.NoError(t, err)
	_TestCommandString(t, cmd, "config")
	require.Equal(t, "config", cmd.CommandName())
	args, is := cmd.Args()
	require.True(t, is)
	require.Equal(t, Cmdline{}, args)	
}

func TestConfigGet (t *testing.T) {
	// config get foo
	cmdName := "config"
	subCmdName := "get"
	subCmdArgs := Cmdline{"foo"}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdArgs)
	cmd := NewConfigCommand(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	subCmd := _TestSubcommandCreation[*ConfigGetCommand](t,
		cmd,
		subCmdName,
		subCmdArgs,
		subCmdString,
	)
	
	require.Equal(t, "foo", subCmd.(*ConfigGetCommand).Key)
}


func TestConfigSet (t *testing.T) {
	// config set foo bar
	cmdName := "config"
	subCmdName := "set"
	subCmdArgs := Cmdline{"foo", "bar"}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdArgs)
	cmd := NewConfigCommand(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	subCmd := _TestSubcommandCreation[*ConfigSetCommand](t,
		cmd,
		subCmdName,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, "foo", subCmd.(*ConfigSetCommand).Key)
	require.Equal(t, "bar", subCmd.(*ConfigSetCommand).Value)
}


func TestConfigShow (t *testing.T) {
	// config set foo bar
	cmdName := "config"
	subCmdName := "show"
	subCmdArgs := Cmdline{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdArgs)
	cmd := NewConfigCommand(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	_TestSubcommandCreation[*ConfigShowCommand](t,
		cmd,
		subCmdName,
		subCmdArgs,
		subCmdString,
	)
}


func TestRunConfigGet (t *testing.T) {
	key := "etcd-domain"
	err := RunConfigGet(key)
	require.Nil(t, err)
}

func TestConfigGetCmd (t *testing.T) {
	dyltPath := GetAndValidateDyltPath(t)
	sCmdline := fmt.Sprintf("%s config get etcd-domain", dyltPath)
	CheckRunCommandSuccess(sCmdline, t)
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


func TestRunConfigShow (t *testing.T) {
	err := RunConfigShow()
	assert.Nil(t, err)
}

func TestConfigShowCmd (t *testing.T) {
	sCmdline := "/Users/chris/src/dylt-dev/dylt/dylt config show"
	CheckRunCommandSuccess(sCmdline, t)
}
