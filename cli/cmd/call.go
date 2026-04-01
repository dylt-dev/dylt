package cmd

import (
	"fmt"
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

	// validate args
	cmdArgs, _ := cmd.Args()
	var v CommandValidator = cmd.CommandValidator()
	if ! v.IsValid(cmdArgs) {
		cmdString, _ := cmd.CommandString()
		errmsg := v.ErrorMessage(cmdArgs)
		return fmt.Errorf("`%s` %s", cmdString, errmsg)
	}

	// init positional params, if any
	if cmd.argmap != nil {
		for i, ptr := range cmd.argmap {
			*ptr = cmdArgs[i]
		}
	}

	return nil
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
