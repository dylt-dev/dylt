package cmd

import (
	"github.com/dylt-dev/dylt/api"
)

type ListCommand struct {
	*BaseCommand
}

func NewListCommand(cmdline Cmdline, parent Command) *ListCommand {
	// list command
	name := "list"
	validator := ArgCountValidator{nExpected: 0}
	cmd := &ListCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_List, nil, validator)}
	cmd.fnRun = func () error { return api.RunList() }
	
	//init flags (if any)
	
	return cmd
}
