package cmd

import (
	"fmt"
	"log/slog"

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
	
	//init flags (if any)
	
	return cmd
}


func (cmd *ListCommand) Run() error {
	slog.Debug("ListCommand.Run()", "args", cmd.Cmdline)

	// parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}
	
	// If help flag set, print usage
	if cmd.Help {
		fmt.Println("halp!")
		cmd.PrintUsage()
		return nil
	}
	
	// execute command
	err = api.RunList()
	if err != nil {
		return err
	}

	return nil
}

