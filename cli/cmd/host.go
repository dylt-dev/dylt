package cmd

import (
	"log/slog"

	"github.com/dylt-dev/dylt/api"
	"github.com/dylt-dev/dylt/common"
)

type HostCommand struct {
	*BaseCommand
}

func NewHostCommand(cmdline Cmdline, parent Command) *HostCommand {
	// host command
	name := "host"
	cmdMap := CommandMap{
		"init": HostInitCommandF.New,
	}
	validator := ArgCountGEValidator{nExpected: 0}
	cmd := &HostCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Host, cmdMap, validator)}
	cmd.isUsageOnNoArgs = true
	
	//init flags (if any)
	
	return cmd
}


func (cmd *HostCommand) Run() error {
	slog.Debug("HostCommand.Run()", "args", cmd.Cmdline)

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

	// Check for 0 args; if so print usage & return
	args, _ := cmd.Args()
	if len(args) == 0 && cmd.UsageOnNoArgs() {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// If CommandMap exists run subcommand
	cmdMap := cmd.CommandMap()
	if cmdMap != nil {
		subCmd, err := cmd.CreateSubCommand()
		if err != nil {
			return err
		}
		err = subCmd.Run()
		return err
	}

	// Execute command
	if cmd.fnRun != nil {
		return cmd.fnRun()
	}

	return nil
}

func RunHost(cmdline Cmdline, parent Command) error {
	slog.Debug("RunHost()", "cmdline", cmdline, "parent", parent)
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

type HostInitCommand struct {
	*BaseCommand
	Gid int // --gid
	Uid int // --uid
}

func NewHostInitCommand(cmdline Cmdline, parent Command) *HostInitCommand {
	// host init command
	name := "host.init"
	validator := ArgCountValidator{nExpected: 0}
	cmd := &HostInitCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Host_Init, nil, validator)}
	cmd.fnRun = func () error { return api.RunHostInit(cmd.Uid, cmd.Gid) }
	
	//init flags (if any)
	cmd.IntVar(&cmd.Gid, "gid", 2000, "gid")
	cmd.IntVar(&cmd.Uid, "uid", 2000, "uid")
	cmd.isUsageOnNoArgs = true

	return cmd
}

func (cmd *HostInitCommand) Run() error {
	slog.Debug("HostInitCommand.Run()", "args", cmd.Cmdline)

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

	// Check for 0 args; if so print usage & return
	args, _ := cmd.Args()
	if len(args) == 0 && cmd.UsageOnNoArgs() {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// If CommandMap exists run subcommand
	cmdMap := cmd.CommandMap()
	if cmdMap != nil {
		subCmd, err := cmd.CreateSubCommand()
		if err != nil {
			return err
		}
		err = subCmd.Run()
		return err
	}

	// Execute command
	if cmd.fnRun != nil {
		return cmd.fnRun()
	}

	return nil
}
