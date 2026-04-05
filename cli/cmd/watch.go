package cmd

import (
	"github.com/dylt-dev/dylt/api"
)

type WatchCommand struct {
	*BaseCommand
}

func NewWatchCommand(cmdline Cmdline, parent Command) *WatchCommand {
	// watch command
	name := "watch"
	cmdMap := CommandMap{
		"script": WatchScriptCommandF.New,
		"svc":    WatchSvcCommandF.New,
	}
	validator := ArgCountGEValidator{nExpected: 0}
	cmd := &WatchCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Watch, cmdMap, validator)}
	cmd.isUsageOnNoArgs = true

	//init flags (if any)

	return cmd
}

// func RunWatch(cmdline Cmdline, parent Command) error {
// 	slog.Debug("RunWatch()", "cmdline", cmdline, "parent", parent)
// 	// Create the subcommand and run it
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

// Usage
//
//	watch script scriptKey targetPath
type WatchScriptCommand struct {
	*BaseCommand
	ScriptKey  string // arg 0
	TargetPath string // arg 1
}

func NewWatchScriptCommand(cmdline Cmdline, parent Command) *WatchScriptCommand {
	// watch script command
	name := "watch.script"
	validator := ArgCountValidator{nExpected: 2}
	cmd := &WatchScriptCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Watch_Script, nil, validator)}
	cmd.argMap = map[int]*string{
		0: &cmd.ScriptKey,
		1: &cmd.TargetPath,
	}
	cmd.fnRun = func() error { return api.RunWatchScript(cmd.ScriptKey, cmd.TargetPath) }

	//init flags (if any)

	return cmd
}

// Usage
//
//	watch svc name
type WatchSvcCommand struct {
	*BaseCommand
	Name string
}

func NewWatchSvcCommand(cmdline Cmdline, parent Command) *WatchSvcCommand {
	// watch svc command
	name := "watch.svc"
	validator := ArgCountValidator{nExpected: 1}
	cmd := &WatchSvcCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Watch_Svc, nil, validator)}
	cmd.argMap = map[int]*string{
		0: &cmd.Name,
	}
	cmd.fnRun = func() error { return api.RunWatchSvc() }

	//init flags (if any)

	return cmd
}
