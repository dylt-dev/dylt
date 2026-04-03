package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVm (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdName := "vm"
	cmdFlags := []string{}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, VmCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString)
	require.IsType(t, &VmCommand{}, cmd)
}


func TestVmHelp (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdName := "vm"
	cmdFlags := []string{"--help"}
	cmdArgs := []string{}
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t, VmCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString).(*VmCommand)
	require.IsType(t, &VmCommand{}, cmd)
	require.True(t, cmd.Help)
}

func TestVmAdd (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	// config get foo
	cmdName := "vm"
	subCmdName := "add"
	name := "name"
	fqdn := "fqdn"
	subCmdFlags := []string{}
	subCmdArgs := []string{name, fqdn}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := VmCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	subCmd := _TestSubcommandCreation[*VmAddCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, fqdn, subCmd.Fqdn)
	require.Equal(t, name, subCmd.Name)
}


func TestVmAddHelp(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdName := "vm"
	subCmdName := "add"
	subCmdFlags := []string{"--help"}
	subCmdArgs := []string{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := VmCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test  subcommand
	subCmd := _TestSubcommandCreation[*VmAddCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.True(t, subCmd.Help)
}

func TestVmAll (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	// config get foo
	cmdName := "vm"
	subCmdName := "all"
	subCmdFlags := []string{}
	subCmdArgs := []string{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := VmCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	_TestSubcommandCreation[*VmAllCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
}

func TestVmDel (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	// config get foo
	cmdName := "vm"
	subCmdName := "del"
	name := "name"
	subCmdFlags := []string{}
	subCmdArgs := []string{name}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := VmCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	subCmd := _TestSubcommandCreation[*VmDelCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, name, subCmd.Name)
}


func TestVmDelHelp(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdName := "vm"
	subCmdName := "del"
	subCmdFlags := []string{"--help"}
	subCmdArgs := []string{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := VmCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test  subcommand
	subCmd := _TestSubcommandCreation[*VmDelCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.True(t, subCmd.Help)
}

func TestVmGet (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	// config get foo
	cmdName := "vm"
	subCmdName := "get"
	name := "name"
	subCmdFlags := []string{}
	subCmdArgs := []string{name}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := VmCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	subCmd := _TestSubcommandCreation[*VmGetCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, name, subCmd.Name)
}

func TestVmList (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	// config get foo
	cmdName := "vm"
	subCmdName := "list"
	subCmdFlags := []string{}
	subCmdArgs := []string{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := VmCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	_TestSubcommandCreation[*VmListCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
}


func TestVmListHelp(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdName := "vm"
	subCmdName := "list"
	subCmdFlags := []string{"--help"}
	subCmdArgs := []string{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := VmCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test  subcommand
	subCmd := _TestSubcommandCreation[*VmListCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.True(t, subCmd.Help)
}


func TestVmSet (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	// config get foo
	cmdName := "vm"
	subCmdName := "set"
	name := "name"
	key := "key"
	value := "value"
	subCmdFlags := []string{}
	subCmdArgs := []string{name, key, value}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := VmCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	subCmd := _TestSubcommandCreation[*VmSetCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, name, subCmd.Name)
	require.Equal(t, key, subCmd.Key)
	require.Equal(t, value, subCmd.Value)
}


func TestVmSetHelp(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdName := "vm"
	subCmdName := "set"
	subCmdFlags := []string{"--help"}
	subCmdArgs := []string{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := VmCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test  subcommand
	subCmd := _TestSubcommandCreation[*VmSetCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.True(t, subCmd.Help)
}