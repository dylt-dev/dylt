package cmd

import (
	"github.com/dylt-dev/dylt/api"
)

type ListOpts struct {
}

func NewListCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	cfg := BaseCommandConfig[ListOpts]{
		name: "list",
		fnRun: func (cmd *BaseCommand[ListOpts]) error { return api.RunList() },
		opts: ListOpts{},
		usage: CreateUsageString(USG_List),
		validator: ArgCountValidator{nExpected: 0},
	}	
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any 

	// subcommand map if any
	
	// done
	return cmd
}
