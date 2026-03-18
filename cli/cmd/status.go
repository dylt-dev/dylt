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
	"flag"
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/lib"
)

type StatusCommand struct {
	*BaseCommand
	Help bool
}

func NewStatusCommand (cmdline Cmdline) *StatusCommand {
	// create command
	flagSet := flag.NewFlagSet("status", flag.ExitOnError)
	cmd := StatusCommand{BaseCommand: &BaseCommand{ Cmdline: cmdline, FlagSet: flagSet}}
	// init flag vars (no flags; nop)
	flagSet.BoolVar(&cmd.Help, "help", false, "blah blah blah")
	return &cmd
}

func (cmd *StatusCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil { return err }
	// validate arg count
	cmdArgs, _ := cmd.Args()
	nExpected := 0
	if len(cmdArgs) != nExpected {
		cmd.PrintUsage()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
			cmd.GetCommandString(),
			nExpected,
			len(cmdArgs))
		}

	return nil
}

func (cmd *StatusCommand) PrintUsage () {
	PrintUsage(USG_Status_Short)
}

func (cmd *StatusCommand) Run() error {
	slog.Debug("StatusCommand.Run()", "args", cmd.Cmdline)
	// parse flags & get positional args
err := cmd.HandleArgs()
	if err != nil { return err }
	// execute command
	// @getit
	err = RunStatus()
	if err != nil { return err }

	return nil
}

func RunStatus() error {
	slog.Debug("RunStatus()")

	lib.RunStatus()

	return nil
}
