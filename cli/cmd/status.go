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
	"github.com/dylt-dev/dylt/api"
)

type StatusOpts struct {
}

func NewStatusCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	name := "status"
	opts := StatusOpts{}
	fnRun := func (cmd *BaseCommand[StatusOpts]) error { return api.RunStatus() }
	cfg := BaseCommandConfig[StatusOpts]{
		name: name,
		fnRun: fnRun,
		opts: opts,
		usage: CreateUsageString(USG_Status),
		validator: ArgCountValidator{nExpected: 0},
	}	
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any
	
	// subcommand map if any
	
	// done
	return cmd
}
