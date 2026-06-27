package cmd

import (
	"github.com/dylt-dev/dylt/api"
)

type WatchOpts struct {
}

func NewWatchCommand(cmdline Cmdline, parent Command) Command {
	// watch command
	cfg := BaseCommandConfig[WatchOpts]{
		name:            "watch",
		opts:            WatchOpts{},
		isUsageOnNoArgs: true,
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
	ScriptKey  string `pos:"0" impl:"api.RunWatchScript"`
	TargetPath string `pos:"1"`
}

func NewWatchScriptCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	cfg := BaseCommandConfig[WatchScriptOpts]{
		name: "watch.script",
		fnRun: func(cmd *BaseCommand[WatchScriptOpts]) error { return api.RunWatchScript(cmd.opts.ScriptKey, cmd.opts.TargetPath) },
		opts: WatchScriptOpts{},
		usage: CreateUsageString(USG_Watch_Script),
		validator: ArgCountValidator{nExpected: 2},
	}	
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any

	// subcommand map if any
	
	// done
	return cmd
}

// Usage
//
//	watch svc name
type WatchSvcOpts struct {
	Name string `pos:"0" impl:"api.RunWatchSvc"`
}

func NewWatchSvcCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	cfg := BaseCommandConfig[WatchSvcOpts]{
		name: "watch.svc",
		fnRun: func(cmd *BaseCommand[WatchSvcOpts]) error { return api.RunWatchSvc() },
		opts: WatchSvcOpts{},
		usage: CreateUsageString(USG_Config_Get),
		validator: ArgCountValidator{nExpected: 1},
	}	
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any

	// subcommand map if any
	
	// done
	return cmd
}
