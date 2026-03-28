package cmd

import (
	"log/slog"

	"github.com/dylt-dev/dylt/common"
	"github.com/dylt-dev/dylt/lib"
)

type CallCommand struct {
	*BaseCommand
	ScriptArgs []string // args 0..n-1
	ScriptPath string   // flag
}

func NewCallCommand(cmdline Cmdline, parent Command) *CallCommand {
	// call command
	name := "call"
	cmd := CallCommand{BaseCommand: NewBaseCommand(name, cmdline, parent)}
	// init flag vars
	cmd.StringVar(&cmd.ScriptPath, "script-path", "/opt/bin/daylight.sh", "script-path")

	return &cmd
}

func (cmd *CallCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}

	// if Help flag is set, no further processing is necessary
	if cmd.Help {
		return nil
	}

	// Check for 0 args; if so print usage & return
	if len(cmd.Cmdline) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// init positional params
	cmd.ScriptArgs = cmd.Cmdline.Args()

	return nil
}

func (cmd *CallCommand) PrintUsage() {
	PrintUsage(USG_Call_Full)
}

func (cmd *CallCommand) Run() error {
	slog.Debug("CallCommand.Run()", "args", cmd.Cmdline)

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

	// Execute command
	err = lib.RunCall(cmd.ScriptPath, cmd.ScriptArgs)
	if err != nil {
		return err
	}

	return nil
}
