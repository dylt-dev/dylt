package cmd

import (
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/eco"
)

type ListCommand struct {
	*BaseCommand
}

func NewListCommand(cmdline Cmdline, parent Command) *ListCommand {
	// list command
	name := "list"
	cmd := &ListCommand{BaseCommand: NewBaseCommand(name, cmdline, parent)}
	
	//init flags (if any)
	
	return cmd
}

func (cmd *ListCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}

	// if Help flag is set, no further processing is necessary
	if cmd.Help {
		return nil
	}

	// validate arg count
	args, _ := cmd.Args()
	nExpected := 0
	if len(args) != nExpected {
		cmd.PrintUsage()
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
			cmdString,
			nExpected,
			len(args))
	}

	// init positional params + subargs (if any)

	return nil
}

func (cmd *ListCommand) PrintUsage() {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_List)
	fmt.Println()
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
	err = RunList()
	if err != nil {
		return err
	}

	return nil
}


func RunList() error {
	// get etcd client + list all keys, one per line
	cli, err := eco.CreateEtcdClientFromConfig()
	if err != nil {
		return err
	}
	kvs, err := cli.List()
	if err != nil {
		return err
	}
	for _, kv := range kvs {
		fmt.Println(string(kv.Key))
	}

	return nil
}
