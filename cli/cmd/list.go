package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	dylt "github.com/dylt-dev/dylt/lib"
)

func CreateListCommand () *cobra.Command {
	command := cobra.Command {
		Use: "list",
		Short: "List all keys",
		Long: "List all keys in etcd cluster",
		RunE: runCommand,
	}
	return &command
}

func runCommand (cmd *cobra.Command, args []string) error {
	cli, err := dylt.NewEtcdClient(viper.GetString("etcd_domain"))
	if err != nil { return err }
	kvs, err := cli.List()
	if err != nil { return err }
	for _, kv := range kvs {
		fmt.Println(string(kv.Key))
	}
	return nil
}