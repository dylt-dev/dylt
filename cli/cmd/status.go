// Copyright © 2025 Chris Lalos <chris@dylt.dev>
// This work is free. You can redistribute it and/or modify it under the
// terms of the Do What The Fig You Want To Public License, Version 2,
// as published by Sam Hocevar and modified by the author.
// See the COPYING file for more details.

// The `status` command provides the status of the local dylt installation. If
// you are wondering how well `dylt` is supported on your local workstation, or
// if it is supported at all, `dylt status` will tell you. Local `dylt`
// installations vary from platform to platform -- `dylt` on Windows works
// differently than `dylt` on OSX, etc -- so the exact output
// of `status` will vary from platform to platform.

package cmd

import (
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/api"
)

type StatusCommand struct {
	*BaseCommand
}

func NewStatusCommand(cmdline Cmdline, parent Command) *StatusCommand {
	// status command
	name := "status"
	validator := ArgCountValidator{nExpected: 0}
	cmd := &StatusCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Status_Short, nil, validator)}
	
	//init flags (if any)

	return cmd
}

func (cmd *StatusCommand) HandleArgs() error {
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

	return nil
}

func (cmd *StatusCommand) Run() error {
	slog.Debug("StatusCommand.Run()", "args", cmd.Cmdline)

	// parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}

	// If help flag set, print usage + return
	if cmd.Help {
		cmd.PrintUsage()
		return nil
	}

	// execute command
	// @getit
	err = api.RunStatus()
	if err != nil {
		return err
	}

	return nil
}