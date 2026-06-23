package cmd

import (
	"github.com/dylt-dev/dylt/lib"
)

type CallOpts struct {
	ScriptPath string `flag:"script-path" default:"/opt/bin/daylight.sh" desc:"script-path" impl:"lib.RunCall"`
}

type CallCommand BaseCommand[CallOpts]

func NewCallCommand(cmdline Cmdline, parent Command) *BaseCommand[CallOpts] {
	// create config object + BaseCommand
	cfg := BaseCommandConfig[CallOpts]{
		name:            "call",
		fnRun:           func(cmd *BaseCommand[CallOpts]) error {
			scriptArgs := cmd.Cmdline.Args()
			err := lib.RunCall(cmd.opts.ScriptPath, scriptArgs)
			return err
		},
		opts:            CallOpts{},
		usage:           CreateUsageString(USG_Call),
		validator:       ArgCountGEValidator{nExpected: 1},
	}
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any

	// subcommand map if any
	
	return cmd
}
