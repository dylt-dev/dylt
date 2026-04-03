package cmd

import (
	"github.com/dylt-dev/dylt/lib"
)

type CallCommandOpts struct {
	ScriptPath string   // flag
}

func NewCallCommand(cmdline Cmdline, parent Command) *BaseCommand {
	// call command
	name := "call"
	validator := ArgCountGEValidator{nExpected: 1}
	opts := CallCommandOpts{}
	cmd := NewBaseCommand(name, cmdline, parent, USG_Call, nil, validator)
	// init flag vars
	cmd.StringVar(&opts.ScriptPath, "script-path", "/opt/bin/daylight.sh", "script-path")
	cmd.fnRun = func () error {
		scriptArgs := cmd.Cmdline.Args()
		err := lib.RunCall(opts.ScriptPath, scriptArgs)
		return err
	}
	cmd.isUsageOnNoArgs = true
	return cmd
}

