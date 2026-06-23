	package cmd

import (
	"github.com/dylt-dev/dylt/lib"
)

type GetOpts struct {
	Key string `pos:"0" impl:"lib.RunGet"`
}

func NewGetCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	cfg := BaseCommandConfig[GetOpts]{
		name:            "get",
		fnRun:           func (cmd *BaseCommand[GetOpts]) error { return lib.RunGet(cmd.opts.Key) },
		opts:            GetOpts{},
		usage:           CreateUsageString(USG_Get),
		validator:       ArgCountValidator{nExpected: 1},
	}
	cmd := NewBaseCommand(cmdline, parent, cfg)
	
	// flags + args if any
	
	// subcommand map if any
	
	// done
	return cmd
}
