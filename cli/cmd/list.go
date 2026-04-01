package cmd

import (
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/api"
	"github.com/dylt-dev/dylt/common"
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
	
	// Check for 0 args; if so print usage & return
	args, _ := cmd.Args()
	if len(args) == 0 && cmd.UsageOnNoArgs() {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// If CommandMap exists run subcommand
	cmdMap := cmd.CommandMap()
	if cmdMap != nil {
		subCmd, err := cmd.CreateSubCommand()
		if err != nil {
			return err
		}
		err = subCmd.Run()
		return err
	}

	// Execute command
	if cmd.fnRun != nil {
		return cmd.fnRun()
	}

	return nil
}
