package cmd

import (
	"github.com/dylt-dev/dylt/api"
)

type MiscOpts struct {
}

func NewMiscCommand(cmdline Cmdline, parent Command) Command {
	// config command
	name := "misc"
	opts := MiscOpts{}
	cfg := BaseCommandConfig[MiscOpts]{
		name:            name,
		isUsageOnNoArgs: true,
		opts:            opts,
		usage:           CreateUsageString(USG_Misc),
		validator:       ArgCountGEValidator{nExpected: 0},
	}
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any

	// subcommand map if any
	cmd.subCommandMap = CommandMap{
		"create-two-node-cluster": CreateTwoNodeClusterCommandF.New,
		"gen-etcd-run-script":     GenEtcdRunScriptCommandF.New,
		"lookup":                  LookupCommandF.New,
	}

	// done
	return cmd
}

// func RunMisc(cmdline Cmdline, parent Command) error {
// 	common.Logger.Debug("RunMisc()", "cmdline", cmdline, "parent", parent)
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

type CreateTwoNodeClusterOpts struct {
}

func NewCreateTwoNodeClusterCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	name := "misc.create-two-node-cluster"
	opts := CreateTwoNodeClusterOpts{}
	fnRun := func(cmd *BaseCommand[CreateTwoNodeClusterOpts]) error { return api.RunCreateTwoNodeCluster() }
	cfg := BaseCommandConfig[CreateTwoNodeClusterOpts]{
		name: name,
		fnRun: fnRun,
		opts: opts,
		usage: CreateUsageString(USG_Misc_TwoNode),
		validator: ArgCountGEValidator{nExpected: 0},
	}	
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// subcommand map if any
	
	// done
	return cmd
}

type GenEtcdRunScriptOpts struct {
}

func NewGenEtcdRunScriptCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	name := "misc.gen-etcd-run-script"
	opts := GenEtcdRunScriptOpts{}
	fnRun := func(cmd *BaseCommand[GenEtcdRunScriptOpts]) error { return api.RunGenEtcdRunScript() }
	cfg := BaseCommandConfig[GenEtcdRunScriptOpts]{
		name: name,
		fnRun: fnRun,
		opts: opts,
		usage: CreateUsageString(USG_Misc_GenScript),
		validator: ArgCountValidator{nExpected: 0},
	}	
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// subcommand map if any
	
	// done
	return cmd
}

type LookupOpts struct {
	Hostname string //arg 0
}

func NewLookupCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	name := "misc.lookup"
	opts := LookupOpts{}
	fnRun := func(cmd *BaseCommand[LookupOpts]) error { return api.RunLookupCommand(cmd.opts.Hostname) }
	cfg := BaseCommandConfig[LookupOpts]{
		name: name,
		fnRun: fnRun,
		opts: opts,
		usage: CreateUsageString(USG_Config_Get),
		validator: ArgCountValidator{nExpected: 1},
	}	
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any
	cmd.argMap = map[int]*string{
		0: &cmd.opts.Hostname,
	}
	// subcommand map if any
	
	// done
	return cmd
}
