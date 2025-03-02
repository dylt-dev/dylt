package cmd

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/lib"
)

type ListCommand struct {
	*flag.FlagSet
}

func NewListCommand () *ListCommand {
	// create command
	flagSet := flag.NewFlagSet("list", flag.PanicOnError)
	cmd := ListCommand { FlagSet: flagSet }
	// init flag vars - (nop - no flags)

	return &cmd
}

func (cmd *ListCommand) HandleArgs(args []string) error {
	// parse flags
	err := cmd.Parse(args)
	if err != nil { return err }
	// validate arg count
	cmdArgs := cmd.Args()
	cmdName := "list"
	nExpected := 0
	if len(cmdArgs) != nExpected { return fmt.Errorf("`%s` expects %d argument(s); received %d", cmdName, nExpected, len(cmdArgs)) }
	// init positional params (nop - no params)

	return nil
}

func (cmd *ListCommand) Run (args []string) error {
	slog.Debug("ListCommand.Run()", "args", args)
	// parse flags & get positional args
	err := cmd.HandleArgs(args)
	if err != nil { return err }
	// execute command
	err = RunList()
	if err != nil { return err }

	return nil
}


func RunList () error {
	// get etcd client + list all keys, one per line
	cli, err := lib.CreateEtcdClientFromConfig()
	if err != nil { return err }
	kvs, err := cli.List()
	if err != nil { return err }
	for _, kv := range kvs {
		fmt.Println(string(kv.Key))
	}
	
	return nil
}