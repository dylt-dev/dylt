package cmd

import (
	"log/slog"

	"github.com/dylt-dev/dylt/api"
	"github.com/dylt-dev/dylt/common"
)

type WatchCommand struct {
	*BaseCommand
}

func NewWatchCommand(cmdline Cmdline, parent Command) *WatchCommand {
	// watch command
	name := "watch"
	cmdMap := CommandMap{
		"script": WatchScriptCommandF.New,
		"svc": WatchSvcCommandF.New,
	}
	validator := ArgCountGEValidator{nExpected: 0}
	cmd := &WatchCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Watch, cmdMap, validator)}

	//init flags (if any)
	
	return cmd
}


func (cmd *WatchCommand) Run() error {
	slog.Debug("WatchCommand.Run()", "args", cmd.Cmdline)

	// Check for 0 args; if so print usage & return
	if len(cmd.Cmdline) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// Parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}

	// If help flag set, print usage + return
	if cmd.Help {
		cmd.PrintUsage()
		return nil
	}

	// If no args, print usage
	args, _ := cmd.Args()
	if len(args) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// if CommandMap exists run subcommand
	cmdMap := cmd.CommandMap()
	if cmdMap != nil {
		subCmd, err := cmd.CreateSubCommand()
		if err != nil {
			return err
		}
		err = subCmd.Run()
		return err
	}

	// execute command
	if cmd.fnRun != nil {
		return cmd.fnRun()
	}

	return nil
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
	cmd.argmap  = map[int]*string {
		0: &cmd.ScriptKey,
		1: &cmd.TargetPath,
	}
	cmd.fnRun = func () error {	return api.RunWatchScript(cmd.ScriptKey, cmd.TargetPath) }
	
	//init flags (if any)
	
	return cmd
}

func (cmd *WatchScriptCommand) Run() error {
	slog.Debug("WatchScriptCommand.Run()", "args", cmd.Cmdline)

	// Parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}

	// if Help flag is set, no further processing is necessary
	if cmd.Help {
		return nil
	}

	// if CommandMap exists run subcommand
	cmdMap := cmd.CommandMap()
	if cmdMap != nil {
		subCmd, err := cmd.CreateSubCommand()
		if err != nil {
			return err
		}
		err = subCmd.Run()
		return err
	}

	// execute command
	if cmd.fnRun != nil {
		return cmd.fnRun()
	}

	return nil
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
	cmd.argmap  = map[int]*string {
		0: &cmd.Name,
	}
	cmd.fnRun = func () error {	return api.RunWatchSvc() }
	
	//init flags (if any)
	
	return cmd
}

func (cmd *WatchSvcCommand) Run() error {
	slog.Debug("WatchSvcCommand.Run()", "args", cmd.Cmdline)

	// Parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}

	// If help flag set, print usage + return
	if cmd.Help {
		cmd.PrintUsage()
		return nil
	}

	// If no args, print usage
	if len(cmd.Cmdline) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// if CommandMap exists run subcommand
	cmdMap := cmd.CommandMap()
	if cmdMap != nil {
		subCmd, err := cmd.CreateSubCommand()
		if err != nil {
			return err
		}
		err = subCmd.Run()
		return err
	}

	// execute command
	if cmd.fnRun != nil {
		return cmd.fnRun()
	}

	return nil
}
