package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	dylt "github.com/dylt-dev/dylt/lib"
)

func CreateGetCommand () *cobra.Command {
	command := cobra.Command {
		Use: "get $key",
		Short: "Get an etcd value by key",
		Long: "Get a value by key from an etcd server. This is simpler than using an etcd client directly.",
		RunE: runGetCommand,
	}
	return &command
}

func runGetCommand (cmd *cobra.Command, args []string) error {
	key := args[0]
	cli, err := dylt.NewEtcdClient("hello.dylt.dev")
	if err != nil { return err }
	val, err := cli.Get(key)
	if err != nil { return err }
	fmt.Println(string(val))
	return nil
}