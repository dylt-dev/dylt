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
	_TestSubcommandCreation[*ConfigShowCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
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


func TestRunConfigShow (t *testing.T) {
	err := RunConfigShow()
	assert.Nil(t, err)
}

func TestConfigShowCmd (t *testing.T) {
	sCmdline := "/Users/chris/src/dylt-dev/dylt/dylt config show"
	lib.CheckRunCommandSuccess(sCmdline, t)
}
