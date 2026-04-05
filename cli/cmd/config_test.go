package cmd

import (
	"fmt"
	"testing"

	"github.com/dylt-dev/dylt/api"
	"github.com/dylt-dev/dylt/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)

	cmdName := "config"
	cmdFlags := []string{}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, ConfigCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString).(*BaseCommand[ConfigOpts])
	require.NotNil(t, cmd)
}

func TestConfigHelp(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)

	cmdName := "list"
	cmdFlags := []string{"--help"}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, ConfigCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString).(*BaseCommand[ConfigOpts])
	require.True(t, cmd.Help())
}

func TestConfigGet(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)

	// config get foo
	cmdName := "config"
	subCmdName := "get"
	subCmdFlags := []string{}
	subCmdArgs := []string{"foo"}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := ConfigCommandF.New(cmdline, nil)
	require.IsType(t, &BaseCommand[ConfigOpts]{}, cmd)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	subCmd := _TestSubcommandCreation[*BaseCommand[ConfigGetOpts]](
		t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, "foo", subCmd.opts.Key)
}

func TestConfigGetHelp(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)

	cmdName := "config"
	subCmdName := "get"
	subCmdFlags := []string{"--help"}
	subCmdArgs := []string{"foo"}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := ConfigCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	subCmd := _TestSubcommandCreation[*BaseCommand[ConfigGetOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.True(t, subCmd.Help())
}

func TestConfigSet(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)

	// config set foo bar
	cmdName := "config"
	subCmdName := "set"
	subCmdFlags := []string{}
	subCmdArgs := Cmdline{"foo", "bar"}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := ConfigCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	subCmd := _TestSubcommandCreation[*BaseCommand[ConfigSetOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, "foo", subCmd.opts.Key)
	require.Equal(t, "bar", subCmd.opts.Value)
}

func TestConfigSetHelp(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)

	cmdName := "config"
	subCmdName := "set"
	subCmdFlags := []string{"--help"}
	subCmdArgs := []string{"foo"}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := ConfigCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	subCmd := _TestSubcommandCreation[*BaseCommand[ConfigSetOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.True(t, subCmd.Help())
}

func TestConfigShow(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)

	// config set foo bar
	cmdName := "config"
	subCmdName := "show"
	subCmdFlags := []string{}
	subCmdArgs := Cmdline{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := ConfigCommandF.New(cmdline, nil)
	cmd.HandleArgs()
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	subCmd := _TestSubcommandCreation[*BaseCommand[ConfigShowOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.False(t, subCmd.Help())
}

func TestConfigShowHelp(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)

	cmdName := "config show"
	cmdFlags := []string{"--help"}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, ConfigShowCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString).(*BaseCommand[ConfigShowOpts])
	require.True(t, cmd.Help())
}

func TestRunConfigGet(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)

	key := "etcd-domain"
	err := api.RunConfigGet(key)
	require.Nil(t, err)
}

func TestConfigGetCmd(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)

	dyltPath := lib.GetAndValidateDyltPath(t)
	sCmdline := fmt.Sprintf("%s config get etcd-domain", dyltPath)
	lib.CheckRunCommandSuccess(sCmdline, t)
}

func TestRunConfigSet(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)

	key := "etcd-domain"
	val := "poo"
	err := api.RunConfigSet(key, val)
	assert.Nil(t, err)
}

func TestConfigSetCmd0(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)

	sCmdline := "/Users/chris/src/dylt-dev/dylt/dylt config set etcd-domain MOO"
	CheckRunCommandSuccessNoOutput(sCmdline, t)
}

func TestConfigSetCmd1(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)

	sCmdline := "/Users/chris/src/dylt-dev/dylt/dylt config set etcd-domain hello.dylt.dev"
	CheckRunCommandSuccessNoOutput(sCmdline, t)
}

// Created this test to figure out why I couldn't create Command objects with
// simple code. Turned out I was forgetting to Parse(). I was used to my
// test_utils helper function doing it for me.
func TestCreateSubcommandRaw(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)

	var err error
	cmdline := Cmdline{"config", "get", "--help"}
	cmd := ConfigCommandF.New(cmdline, nil)
	err = cmd.Parse()
	require.NoError(t, err)
	t.Logf("cmd=%v", cmd)
	t.Logf("cmd.CommandName()=%v", cmd.CommandName())
	require.NotNil(t, cmd)
	subCmd, err := cmd.CreateSubCommand()
	err = subCmd.Parse()
	require.NoError(t, err)
	require.NotNil(t, subCmd)
	require.True(t, subCmd.(*BaseCommand[ConfigGetOpts]).Help())
}

func TestRunConfigShow(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)

	err := api.RunConfigShow()
	assert.Nil(t, err)
}

func TestConfigShowCmd(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)

	sCmdline := "/Users/chris/src/dylt-dev/dylt/dylt config show"
	lib.CheckRunCommandSuccess(sCmdline, t)
}
