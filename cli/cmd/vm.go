package cmd

import (
	// "encoding/json"

	"github.com/dylt-dev/dylt/api"
)

type VmOpts struct {
}

func NewVmCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	cfg := BaseCommandConfig[VmOpts]{
		name:            "vm",
		opts:            VmOpts{},
		isUsageOnNoArgs: true,
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
	Name string `pos:"0" impl:"api.RunVmAdd"`
	Fqdn string `pos:"1"`
}

func NewVmAddCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	cfg := BaseCommandConfig[VmAddOpts]{
		name: "vm.add",
		fnRun: func(cmd *BaseCommand[VmAddOpts]) error { return api.RunVmAdd(cmd.opts.Name, cmd.opts.Fqdn) },
		opts: VmAddOpts{},
		usage:           CreateUsageString(USG_Vm_Add),
		validator: ArgCountValidator{nExpected: 2},
	}	
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any

	// subcommand map if any
	
	// done
	return cmd
}

type VmAllOpts struct {
}

func NewVmAllCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	cfg := BaseCommandConfig[VmAllOpts]{
		name: "vm.all",
		fnRun: func(cmd *BaseCommand[VmAllOpts]) error { return api.RunVmAll() },
		opts: VmAllOpts{},
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
	Name string `pos:"0" impl:"api.RunVmDel"`
}

func NewVmDelCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	cfg := BaseCommandConfig[VmDelOpts]{
		name: "vm.del",
		fnRun: func(cmd *BaseCommand[VmDelOpts]) error { return api.RunVmDel(cmd.opts.Name) },
		opts: VmDelOpts{},
		usage: CreateUsageString(USG_Vm_Del),
		validator: ArgCountValidator{nExpected: 1},
	}	
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any

	// subcommand map if any
	
	// done
	return cmd
}

// Usage
//
//	vm get vmName
type VmGetOpts struct {
	Name string `pos:"0" impl:"api.RunVmGet"`
}

func NewVmGetCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	cfg := BaseCommandConfig[VmGetOpts]{
		name: "vm.get",
		fnRun: func(cmd *BaseCommand[VmGetOpts]) error { return api.RunVmGet(cmd.opts.Name) },
		opts: VmGetOpts{},
		usage: CreateUsageString(USG_Vm_Get),
		validator: ArgCountValidator{nExpected: 1},
	}	
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any

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
	cfg := BaseCommandConfig[VmListOpts]{
		name: "vm.list",
		fnRun: func(cmd *BaseCommand[VmListOpts]) error { return api.RunVmList() },
		opts: VmListOpts{},
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
	Name  string `pos:"0" impl:"api.RunVmSet"`
	Key   string `pos:"1"`
	Value string `pos:"2"`
}

func NewVmSetCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	cfg := BaseCommandConfig[VmSetOpts]{
		name: "vm.set",
		fnRun: func(cmd *BaseCommand[VmSetOpts]) error { return api.RunVmSet(cmd.opts.Name, cmd.opts.Key, cmd.opts.Value) },
		opts: VmSetOpts{},
		usage: CreateUsageString(USG_Vm_Set),
		validator: ArgCountValidator{nExpected: 3},
	}	
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any

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
