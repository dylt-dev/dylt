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
		Args: cobra.ExactArgs(1),
	}
	command.Flags().Bool("keys", false, "--keys")
	
	return &command
}

func runGetCommand (cmd *cobra.Command, args []string) error {
	key := args[0]
	cli, err := dylt.NewEtcdClient("hello.dylt.dev")
	if err != nil { return err }
	flKeys, err := cmd.Flags().GetBool("keys")
	if err != nil { return err }
	if flKeys {
		kvs, err := cli.GetKeys(key)
		if err != nil { return err }
		for _, kv := range kvs {
			fmt.Println(kv)
		}
	} else {
		fmt.Printf("flKeys=%t\n", flKeys)
		val, err := cli.Get(key)
		if err != nil { return err }
		fmt.Println(string(val))
	}
	return nil
}