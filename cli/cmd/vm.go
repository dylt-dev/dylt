package cmd

import (
	// "encoding/json"

	"github.com/dylt-dev/dylt/api"
)

type VmOpts struct {
}

func NewVmCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	name := "vm"
	opts := VmOpts{}
	cfg := BaseCommandConfig[VmOpts]{
		name:            name,
		isUsageOnNoArgs: true,
		opts:            opts,
		usage:           CreateUsageString(USG_Vm),
		validator: ArgCountGEValidator{nExpected: 0},
	}
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any

	// subcommand map if any
	cmd.subCommandMap = CommandMap{
		"add":  VmAddCommandF.New,
		"all":  VmAllCommandF.New,
		"del":  VmDelCommandF.New,
		"get":  VmGetCommandF.New,
		"list": VmListCommandF.New,
		"set":  VmSetCommandF.New,
	}

	// done
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

type VmAddOpts struct {
	Name string // arg 0
	Fqdn string // arg 1
}

func NewVmAddCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	name := "vm.add"
	opts := VmAddOpts{}
	fnRun := func(cmd *BaseCommand[VmAddOpts]) error { return api.RunVmAdd(cmd.opts.Name, cmd.opts.Fqdn) }
	cfg := BaseCommandConfig[VmAddOpts]{
		name: name,
		fnRun: fnRun,
		opts: opts,
		usage:           CreateUsageString(USG_Vm_Add),
		validator: ArgCountValidator{nExpected: 2},
	}	
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any
	cmd.argMap = map[int]*string{
		0: &cmd.opts.Name,
		1: &cmd.opts.Fqdn,
	}
	// subcommand map if any
	
	// done
	return cmd
}

type VmAllOpts struct {
}

func NewVmAllCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	name := "vm.all"
	opts := VmAllOpts{}
	fnRun := func(cmd *BaseCommand[VmAllOpts]) error { return api.RunVmAll() }
	cfg := BaseCommandConfig[VmAllOpts]{
		name: name,
		fnRun: fnRun,
		opts: opts,
		usage: CreateUsageString(USG_Vm_All),
		validator: ArgCountValidator{nExpected: 0},
	}	
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any
	
	// subcommand map if any
	
	// done
	return cmd
}

// Usage
//
//	vm del vmName
type VmDelOpts struct {
	Name string // arg 0
}

func NewVmDelCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	name := "vm.del"
	opts := VmDelOpts{}
	fnRun := func(cmd *BaseCommand[VmDelOpts]) error { return api.RunVmDel(cmd.opts.Name) }
	cfg := BaseCommandConfig[VmDelOpts]{
		name: name,
		fnRun: fnRun,
		opts: opts,
		usage: CreateUsageString(USG_Vm_Del),
		validator: ArgCountValidator{nExpected: 1},
	}	
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any
	cmd.argMap = map[int]*string{
		0: &cmd.opts.Name,
	}
	
	// subcommand map if any
	
	// done
	return cmd
}

// Usage
//
//	vm get vmName
type VmGetOpts struct {
	Name string // arg 0
}

func NewVmGetCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	name := "vm.get"
	opts := VmGetOpts{}
	fnRun := func(cmd *BaseCommand[VmGetOpts]) error { return api.RunVmGet(cmd.opts.Name) }
	cfg := BaseCommandConfig[VmGetOpts]{
		name: name,
		fnRun: fnRun,
		opts: opts,
		usage: CreateUsageString(USG_Vm_Get),
		validator: ArgCountValidator{nExpected: 1},
	}	
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any
	cmd.argMap = map[int]*string{
		0: &cmd.opts.Name,
	}
	
	// subcommand map if any
	
	// done
	return cmd
}

// Usage
//
//	vm list
type VmListOpts struct {
}

func NewVmListCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	name := "vm.list"
	opts := VmListOpts{}
	fnRun := func(cmd *BaseCommand[VmListOpts]) error { return api.RunVmList() }
	cfg := BaseCommandConfig[VmListOpts]{
		name: name,
		fnRun: fnRun,
		opts: opts,
		usage: CreateUsageString(USG_Vm_List),
		validator: ArgCountValidator{nExpected: 0},
	}	
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any
	
	// subcommand map if any
	
	// done
	return cmd
}

// Usage
//
//	vm set vmName key val
type VmSetOpts struct {
	Name  string // arg 0
	Key   string // arg 1
	Value string // arg 2
}

func NewVmSetCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	name := "vm.set"
	opts := VmSetOpts{}
	fnRun := func(cmd *BaseCommand[VmSetOpts]) error { return api.RunVmSet(cmd.opts.Name, cmd.opts.Key, cmd.opts.Value) }
	cfg := BaseCommandConfig[VmSetOpts]{
		name: name,
		fnRun: fnRun,
		opts: opts,
		usage: CreateUsageString(USG_Vm_Set),
		validator: ArgCountValidator{nExpected: 3},
	}	
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any
	cmd.argMap = ArgMap{
		0: &cmd.opts.Name,
		1: &cmd.opts.Key,
		2: &cmd.opts.Value,
	}
	
	// subcommand map if any
	
	// done
	return cmd
}

// type VmShowCommand struct {
// 	*flag.FlagSet
// }

// New\w\+CommandFunc NewVmShowCommand () Command {
// 	// create command
// 	flagSet := flag.NewFlagSet("vm.show", flag.ExitOnError)
// 	cmd := VmShowCommand{FlagSet: flagSet}
// 	// init flag vars (nop -- no flags)

// 	return &cmd
// }
