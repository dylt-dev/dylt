package cmd

import (
	"github.com/dylt-dev/dylt/lib"
)

type CallCommand struct {
	*BaseCommand
	ScriptPath string   // flag
}

func NewCallCommand(cmdline Cmdline, parent Command) *CallCommand {
	// call command
	name := "call"
	validator := ArgCountGEValidator{nExpected: 1}
	cmd := CallCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Call, nil, validator)}
	// init flag vars
	cmd.StringVar(&cmd.ScriptPath, "script-path", "/opt/bin/daylight.sh", "script-path")
	cmd.fnRun = func () error {
		scriptArgs := cmd.Cmdline.Args()
		err := lib.RunCall(cmd.ScriptPath, scriptArgs)
		return err
	}
	cmd.isUsageOnNoArgs = true
	return &cmd
}

