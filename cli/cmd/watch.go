package cmd

import (
	"fmt"
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
	cmd := &WatchCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Watch, cmdMap)}
	
	//init flags (if any)
	
	return cmd
}


func (cmd *WatchCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}

	// if Help flag is set, no further processing is necessary
	if cmd.Help {
		return nil
	}

	// validate arg count
	cmdArgs, _ := cmd.Args()
	nExpected := 0
	if len(cmdArgs) < nExpected {
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("`%s` expects >=%d argument(s); received %d",
			cmdString,
			nExpected,
			len(cmdArgs))
	}

	// init positional params (nop - no params)

	return nil
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

	// Execute command
	err = RunWatch(args, cmd)
	return err
}

func RunWatch(cmdline Cmdline, parent Command) error {
	slog.Debug("RunWatch()", "cmdline", cmdline, "parent", parent)
	// Create the subcommand and run it
	subCmd, err := parent.CreateSubCommand()
	if err != nil {
		return err
	}
	err = subCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

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
	cmd := &WatchScriptCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Watch_Script, nil)}
	
	//init flags (if any)
	
	return cmd
}

func (cmd *WatchScriptCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}

	// if Help flag is set, no further processing is necessary
	if cmd.Help {
		return nil
	}

	// validate arg count
	cmdArgs, _ := cmd.Args()
	nExpected := 2
	if len(cmdArgs) != nExpected {
		cmd.PrintUsage()
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
			cmdString,
			nExpected,
			len(cmd.Cmdline))
	}

	// init positional params
	cmd.ScriptKey = cmdArgs[0]
	cmd.TargetPath = cmdArgs[1]

	return nil
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

	// Execute command
	err = api.RunWatchScript(cmd.ScriptKey, cmd.TargetPath)
	return err
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
	cmd := &WatchSvcCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Watch_Svc, nil)}
	
	//init flags (if any)
	
	return cmd
}

func (cmd *WatchSvcCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}

	// if Help flag is set, no further processing is necessary
	if cmd.Help {
		return nil
	}

	// validate arg count
	cmdArgs, _ := cmd.Args()
	nExpected := 1
	if len(cmdArgs) != nExpected {
		cmd.PrintUsage()
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
			cmdString,
			nExpected,
			len(cmdArgs))
	}

	// init positional params
	cmd.Name = cmdArgs[0]

	return nil
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

	// Execute command
	err = api.RunWatchSvc()
	return err

}
