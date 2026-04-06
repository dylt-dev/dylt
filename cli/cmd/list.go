package cmd

import (
	"github.com/dylt-dev/dylt/api"
)

type ListOpts struct {
}

func NewListCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	name := "list"
	opts := ListOpts{}
	fnRun := func (cmd *BaseCommand[ListOpts]) error { return api.RunList() }
	cfg := BaseCommandConfig[ListOpts]{
		name: name,
		fnRun: fnRun,
		opts: opts,
		usage: CreateUsageString(USG_Config_Get),
		validator: ArgCountValidator{nExpected: 0},
	}	
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any 

	// subcommand map if any
	
	// done
	return cmd
}
