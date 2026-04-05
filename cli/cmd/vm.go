package cmd

import (
	// "encoding/json"

	"github.com/dylt-dev/dylt/api"
)

type VmCommand struct {
	*BaseCommand
}

func NewVmCommand(cmdline Cmdline, parent Command) *VmCommand {
	// vm command
	name := "vm"
	cmdMap := CommandMap{
		"add":  VmAddCommandF.New,
		"all":  VmAllCommandF.New,
		"del":  VmDelCommandF.New,
		"get":  VmGetCommandF.New,
		"list": VmListCommandF.New,
		"set":  VmSetCommandF.New,
	}
	validator := ArgCountGEValidator{nExpected: 0}
	cmd := &VmCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Vm, cmdMap, validator)}
	cmd.isUsageOnNoArgs = true

	//init flags (if any)

	return cmd
}

// func RunVm(cmdline Cmdline, parent *VmCommand) error {
// 	slog.Debug("RunVm()", "cmdline", cmdline, "parent", parent)
// 	// create the subcommand and run it
// 	subCmd, err := parent.CreateSubCommand()
// 	if err != nil {
// 		return err
// 	}
// 	err = subCmd.Run()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

type VmAddCommand struct {
	*BaseCommand
	Name string // arg 0
	Fqdn string // arg 1
}

func NewVmAddCommand(cmdline Cmdline, parent Command) *VmAddCommand {
	// vm add command
	name := "vm.add"
	validator := ArgCountValidator{nExpected: 2}
	cmd := &VmAddCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Vm_Add, nil, validator)}
	cmd.argMap = map[int]*string{
		0: &cmd.Name,
		1: &cmd.Fqdn,
	}
	cmd.fnRun = func() error { return api.RunVmAdd(cmd.Name, cmd.Fqdn) }

	//init flags (if any)

	return cmd

}

type VmAllCommand struct {
	*BaseCommand
}

func NewVmAllCommand(cmdline Cmdline, parent Command) *VmAllCommand {
	// vm all command
	name := "vm.all"
	validator := ArgCountValidator{nExpected: 0}
	cmd := &VmAllCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Vm_All, nil, validator)}
	cmd.fnRun = func() error { return api.RunVmAll() }

	// init flags (if any)

	return cmd
}

// Usage
//
//	vm del vmName
type VmDelCommand struct {
	*BaseCommand
	Name string // arg 0
}

func NewVmDelCommand(cmdline Cmdline, parent Command) *VmDelCommand {
	// vm del command
	name := "vm.del"
	validator := ArgCountValidator{nExpected: 1}
	cmd := &VmDelCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Vm_Del, nil, validator)}
	cmd.argMap = map[int]*string{
		0: &cmd.Name,
	}
	cmd.fnRun = func() error { return api.RunVmDel(cmd.Name) }

	return cmd
}

// Usage
//
//	vm get vmName
type VmGetCommand struct {
	*BaseCommand
	Name string // arg 0
}

func NewVmGetCommand(cmdline Cmdline, parent Command) *VmGetCommand {
	// vm get command
	name := "vm.get"
	validator := ArgCountValidator{nExpected: 1}
	cmd := &VmGetCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Vm_Get, nil, validator)}
	cmd.argMap = map[int]*string{
		0: &cmd.Name,
	}
	cmd.fnRun = func() error { return api.RunVmGet(cmd.Name) }

	//init flags (if any)

	return cmd
}

// Usage
//
//	vm list
type VmListCommand struct {
	*BaseCommand
}

func NewVmListCommand(cmdline Cmdline, parent Command) *VmListCommand {
	// vm list command
	name := "vm.list"
	validator := ArgCountValidator{nExpected: 0}
	cmd := &VmListCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Vm_List, nil, validator)}
	cmd.fnRun = func() error { return api.RunVmList() }

	//init flags (if any)

	return cmd
}

// Usage
//
//	vm set vmName key val
type VmSetCommand struct {
	*BaseCommand
	Name  string // arg 0
	Key   string // arg 1
	Value string // arg 2
}

func NewVmSetCommand(cmdline Cmdline, parent Command) *VmSetCommand {
	// vm set command
	name := "vm.set"
	validator := ArgCountValidator{nExpected: 3}
	cmd := &VmSetCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Vm_Set, nil, validator)}
	cmd.argMap = map[int]*string{
		0: &cmd.Name,
		1: &cmd.Key,
		2: &cmd.Value,
	}
	cmd.fnRun = func() error { return api.RunVmSet(cmd.Name, cmd.Key, cmd.Value) }

	//init flags (if any)

	return cmd
}

// type VmShowCommand struct {
// 	*flag.FlagSet
// }

// func NewVmShowCommand () *VmShowCommand {
// 	// create command
// 	flagSet := flag.NewFlagSet("vm.show", flag.ExitOnError)
// 	cmd := VmShowCommand{FlagSet: flagSet}
// 	// init flag vars (nop -- no flags)

// 	return &cmd
// }
