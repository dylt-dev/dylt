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
	flagSet := flag.NewFlagSet("list", flag.PanicOnError)
	return &ListCommand {
		FlagSet: flagSet,
	}
}

func (cmd *ListCommand) Run (args []string) error {
	slog.Debug("ListCommand.Run()", "args", args)
	err := cmd.Parse(args)
	if err != nil { return err }
	cli, err := lib.CreateEtcdClientFromConfig()
	if err != nil { return err }
	kvs, err := cli.List()
	if err != nil { return err }
	for _, kv := range kvs {
		fmt.Println(string(kv.Key))
	}
	return nil
}