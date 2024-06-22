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
	arg := args[0]
	cli, err := dylt.NewEtcdClient("hello.dylt.dev")
	if err != nil { return err }
	flKeys, err := cmd.Flags().GetBool("keys")
	if err != nil { return err }
	if flKeys {
		prefix := arg
		return getKeys(cli, prefix)
	} else {
		key := arg
		val, err := cli.Get(key)
		if err != nil { return err }
		if val == nil { return nil }
		fmt.Println(string(val))
	}
	return nil
}


func getKeys (cli *dylt.EtcdClient, prefix string) error {
	kvs, err := cli.GetKeys(prefix)
	if err != nil { return err }
	for _, kv := range kvs {
		fmt.Println(kv)
	}
	return nil
}