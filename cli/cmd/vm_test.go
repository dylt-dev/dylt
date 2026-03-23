package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVm (t *testing.T) {
	cmdName := "vm"
	cmdFlags := []string{}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, NewVmCommand, cmdName, cmdFlags, cmdArgs, cmdString)
	require.IsType(t, &VmCommand{}, cmd)
}

func TestVmAdd (t *testing.T) {
	// config get foo
	cmdName := "vm"
	subCmdName := "add"
	name := "name"
	fqdn := "fqdn"
	subCmdFlags := []string{}
	subCmdArgs := []string{name, fqdn}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := NewVmCommand(cmdline, nil)
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

func TestVmAll (t *testing.T) {
	// config get foo
	cmdName := "vm"
	subCmdName := "all"
	subCmdFlags := []string{}
	subCmdArgs := []string{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := NewVmCommand(cmdline, nil)
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
	// config get foo
	cmdName := "vm"
	subCmdName := "del"
	name := "name"
	subCmdFlags := []string{}
	subCmdArgs := []string{name}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := NewVmCommand(cmdline, nil)
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

func TestVmGet (t *testing.T) {
	// config get foo
	cmdName := "vm"
	subCmdName := "get"
	name := "name"
	subCmdFlags := []string{}
	subCmdArgs := []string{name}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := NewVmCommand(cmdline, nil)
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
	// config get foo
	cmdName := "vm"
	subCmdName := "list"
	subCmdFlags := []string{}
	subCmdArgs := []string{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := NewVmCommand(cmdline, nil)
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

func TestVmSet (t *testing.T) {
	// config get foo
	cmdName := "vm"
	subCmdName := "set"
	name := "name"
	key := "key"
	value := "value"
	subCmdFlags := []string{}
	subCmdArgs := []string{name, key, value}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := NewVmCommand(cmdline, nil)
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