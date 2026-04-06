package cmd

import (
	"fmt"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

func TestVm (t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	cmdName := "vm"
	cmdFlags := []string{}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, VmCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString)
	require.IsType(t, &BaseCommand[VmOpts]{}, cmd)
}


func TestVmHelp (t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	cmdName := "vm"
	cmdFlags := []string{"--help"}
	cmdArgs := []string{}
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t, VmCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString).(*BaseCommand[VmOpts])
	require.IsType(t, &BaseCommand[VmOpts]{}, cmd)
	require.True(t, cmd.Help())
}

func TestVmAdd (t *testing.T) {
	fnTeardown := common.Setup(t)
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
	subCmd := _TestSubcommandCreation[*BaseCommand[VmAddOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, fqdn, subCmd.opts.Fqdn)
	require.Equal(t, name, subCmd.opts.Name)
}


func TestVmAddHelp(t *testing.T) {
	fnTeardown := common.Setup(t)
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
	subCmd := _TestSubcommandCreation[*BaseCommand[VmAddOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.True(t, subCmd.Help())
}

func TestVmAll (t *testing.T) {
	fnTeardown := common.Setup(t)
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
	_TestSubcommandCreation[*BaseCommand[VmAllOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
}

func TestVmDel (t *testing.T) {
	fnTeardown := common.Setup(t)
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
	subCmd := _TestSubcommandCreation[*BaseCommand[VmDelOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, name, subCmd.opts.Name)
}


func TestVmDelHelp(t *testing.T) {
	fnTeardown := common.Setup(t)
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
	subCmd := _TestSubcommandCreation[*BaseCommand[VmDelOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.True(t, subCmd.Help())
}

func TestVmGet (t *testing.T) {
	fnTeardown := common.Setup(t)
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
	subCmd := _TestSubcommandCreation[*BaseCommand[VmGetOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, name, subCmd.opts.Name)
}

func TestVmList (t *testing.T) {
	fnTeardown := common.Setup(t)
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
	_TestSubcommandCreation[*BaseCommand[VmListOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
}


func TestVmListHelp(t *testing.T) {
	fnTeardown := common.Setup(t)
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
	subCmd := _TestSubcommandCreation[*BaseCommand[VmListOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.True(t, subCmd.Help())
}


func TestVmSet (t *testing.T) {
	fnTeardown := common.Setup(t)
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
	subCmd := _TestSubcommandCreation[*BaseCommand[VmSetOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, name, subCmd.opts.Name)
	require.Equal(t, key, subCmd.opts.Key)
	require.Equal(t, value, subCmd.opts.Value)
}


func TestVmSetHelp(t *testing.T) {
	fnTeardown := common.Setup(t)
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
	subCmd := _TestSubcommandCreation[*BaseCommand[VmSetOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.True(t, subCmd.Help())
}