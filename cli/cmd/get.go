	package cmd

import (
	"github.com/dylt-dev/dylt/lib"
)

type GetOpts struct {
	Key string // arg 0
}

func NewGetCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	fnRun := func (cmd *BaseCommand[GetOpts]) error { return lib.RunGet(cmd.opts.Key) }

	cfg := BaseCommandConfig[GetOpts]{
		name:            "get",
		fnRun:           fnRun,
		opts:            GetOpts{},
		usage:           CreateUsageString(USG_Get),
		validator:       ArgCountValidator{nExpected: 1},
	}
	cmd := NewBaseCommand(cmdline, parent, cfg)
	
	// flags + args if any 
	cmd.argMap = ArgMap{
		0: &cmd.opts.Key,
	}
	
	// subcommand map if any
	
	// done
	return cmd
}
