package cmd

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/eco"
)

type ListCommand struct {
	*BaseCommand
}

func NewListCommand (cmdline Cmdline) *ListCommand {
	// create command
	flagSet := flag.NewFlagSet("list", flag.PanicOnError)
	cmd := ListCommand { BaseCommand: &BaseCommand{Cmdline: cmdline, FlagSet: flagSet} }
	// init flag vars - (nop - no flags)

	return &cmd
}

func (cmd *ListCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil { return err }
	// validate arg count
	cmdArgs, _ := cmd.Args()
	nExpected := 0
	if len(cmd.Cmdline) != nExpected {
		cmd.PrintUsage()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
		                  cmd.GetCommandString(),
						  nExpected,
						  len(cmdArgs))
	}
	// init positional params + subargs

	return nil
}

func (cmd *ListCommand) PrintUsage () {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_List)
	fmt.Println()
}

func (cmd *ListCommand) Run () error {
	slog.Debug("ListCommand.Run()", "args", cmd.Cmdline)
	// parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil { return err }
	// execute command
	err = RunList()
	if err != nil { return err }

	return nil
}


func RunList () error {
	// get etcd client + list all keys, one per line
	cli, err := eco.CreateEtcdClientFromConfig()
	if err != nil { return err }
	kvs, err := cli.List()
	if err != nil { return err }
	for _, kv := range kvs {
		fmt.Println(string(kv.Key))
	}
	
	return nil
}