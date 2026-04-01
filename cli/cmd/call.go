package cmd

import (
	"log/slog"

	"github.com/dylt-dev/dylt/common"
	"github.com/dylt-dev/dylt/lib"
)

type CallCommand struct {
	*BaseCommand
	ScriptPath string   // flag
}

func NewCallCommand(cmdline Cmdline, parent Command) *CallCommand {
	// call command
	name := "call"
	validator := ArgCountGEValidator{nExpected: 1}
	cmd := CallCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Call, nil, validator)}
	// init flag vars
	cmd.StringVar(&cmd.ScriptPath, "script-path", "/opt/bin/daylight.sh", "script-path")

	return &cmd
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

	// Check for 0 args; if so print usage & return
	if len(cmd.Cmdline) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}
	// Execute command
	// init positional params
	scriptArgs := cmd.Cmdline.Args()
	err = lib.RunCall(cmd.ScriptPath, scriptArgs)
	if err != nil {
		return err
	}

	return nil
}
