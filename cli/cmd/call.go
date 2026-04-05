package cmd

import (
	"github.com/dylt-dev/dylt/lib"
)

type CallOpts struct {
	ScriptPath string // flag
}

type CallCommand BaseCommand[CallOpts]

func NewCallCommand(cmdline Cmdline, parent Command) *BaseCommand[CallOpts] {
	// create config object + BaseCommand
	name := "call"
	opts := CallOpts{}
	fnRun := func(cmd *BaseCommand[CallOpts]) error {
		scriptArgs := cmd.Cmdline.Args()
		err := lib.RunCall(opts.ScriptPath, scriptArgs)
		return err
	}
	cfg := BaseCommandConfig[CallOpts]{
		name:            name,
		fnRun:           fnRun,
		opts:            opts,
		usage:           CreateUsageString(USG_Call),
		validator:       ArgCountGEValidator{nExpected: 1},
	}
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// Add flags & args
	cmd.StringVar(&cmd.opts.ScriptPath, "script-path", "/opt/bin/daylight.sh", "script-path")

	// subcommand map if any
	
	return cmd
}
