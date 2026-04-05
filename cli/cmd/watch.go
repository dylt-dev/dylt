package cmd

import (
	"github.com/dylt-dev/dylt/api"
)

type WatchOpts struct {
}

func NewWatchCommand(cmdline Cmdline, parent Command) Command {
	// config command
	name := "watch"
	opts := WatchOpts{}
	cfg := BaseCommandConfig[WatchOpts]{
		name:            name,
		isUsageOnNoArgs: true,
		opts:            opts,
		usage:           CreateUsageString(USG_Watch),
		validator: ArgCountGEValidator{nExpected: 0},
	}
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any

	// subcommand map if any
	cmd.subCommandMap = CommandMap{
		"script": WatchScriptCommandF.New,
		"svc":    WatchSvcCommandF.New,
	}

	// done
	return cmd
}

// func RunWatch(cmdline Cmdline, parent Command) error {
// 	slog.Debug("RunWatch()", "cmdline", cmdline, "parent", parent)
// 	// Create the subcommand and run it
// 	subCmd, err := parent.CreateSubCommand()
// 	if err != nil {

// 	}
// 	err = subCmd.Run()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// Usage
//
//	watch script scriptKey targetPath
type WatchScriptOpts struct {
	ScriptKey  string // arg 0
	TargetPath string // arg 1
}

func NewWatchScriptCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	name := "watch.script"
	opts := WatchScriptOpts{}
	fnRun := func(cmd *BaseCommand[WatchScriptOpts]) error { return api.RunWatchScript(cmd.opts.ScriptKey, cmd.opts.TargetPath) }
	cfg := BaseCommandConfig[WatchScriptOpts]{
		name: name,
		fnRun: fnRun,
		opts: opts,
		usage: CreateUsageString(USG_Watch_Script),
		validator: ArgCountValidator{nExpected: 2},
	}	
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any
	cmd.argMap = map[int]*string{
		0: &cmd.opts.ScriptKey,
		1: &cmd.opts.TargetPath,
	}
	
	// subcommand map if any
	
	// done
	return cmd
}

// Usage
//
//	watch svc name
type WatchSvcOpts struct {
	Name string
}

func NewWatchSvcCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	name := "watch.svc"
	opts := WatchSvcOpts{}
	fnRun := func(cmd *BaseCommand[WatchSvcOpts]) error { return api.RunWatchSvc() }
	cfg := BaseCommandConfig[WatchSvcOpts]{
		name: name,
		fnRun: fnRun,
		opts: opts,
		usage: CreateUsageString(USG_Config_Get),
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
