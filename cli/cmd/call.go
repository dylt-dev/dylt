package cmd

import (
	"github.com/dylt-dev/dylt/lib"
)

type CallOpts struct {
	ScriptPath string // flag
}

type CallCommand BaseCommand[CallOpts]

func NewCallCommand(cmdline Cmdline, parent Command) *BaseCommand[CallOpts] {
	// call command
	name := "call"
	opts := CallOpts{}
	fnRun := func(cmd *BaseCommand[CallOpts]) error {
		scriptArgs := cmd.Cmdline.Args()
		err := lib.RunCall(opts.ScriptPath, scriptArgs)
		return err
	}

	// create config object + BaseCommand
	cfg := BaseCommandConfig[CallOpts]{
		name:            name,
		opts:            opts
		validator:       ArgCountGEValidator{nExpected: 1},
		isUsageOnNoArgs: true,
		fnRun:           fnRun,
		usage:           CreateUsageString(USG_Call),
	}
	cmd := NewBaseCommand[CallOpts](cmdline, parent, cfg)

	// Add flags & args
	cmd.StringVar(&opts.ScriptPath, "script-path", "/opt/bin/daylight.sh", "script-path")

	return cmd
}
