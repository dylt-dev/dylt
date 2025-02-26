package cmd

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/lib"
)

type GetCommand struct {
	*flag.FlagSet
}

func NewGetCommand () *GetCommand {
	flagSet := flag.NewFlagSet("get", flag.PanicOnError)
	return &GetCommand {
		FlagSet: flagSet,
	}
}

func (cmd *GetCommand) Run (args []string) error {
	slog.Debug("GetCommand.Run()", "args", args)
	err := cmd.Parse(args)
	if err != nil { return err }
	cmdArgs := cmd.Args()
	if len(cmdArgs) != 1 {
		return fmt.Errorf("get takes 1 argument; got %d", len(cmdArgs))
	}
	key := cmdArgs[0]
	cli, err := lib.CreateEtcdClientFromConfig()
	if err != nil { return err }
	val, err := cli.Get(key)
	if err != nil { return err }
	fmt.Printf("\n%s\n", val)
	return nil
}

// func GetCommand () *cobra.Command {
// 	command := cobra.Command {
// 		Use: "get $key",
// 		Short: "Get an etcd value by key",
// 		Long: "Get a value by key from an etcd server. This is simpler than using an etcd client directly.",
// 		RunE: runGetCommand,
// 		Args: cobra.ExactArgs(1),
// 	}
// 	command.Flags().Bool("keys", false, "--keys")
	
// 	return &command
// }

// func runGetCommand (cmd *cobra.Command, args []string) error {
// 	arg := args[0]
// 	cli, err := dylt.CreateEtcdClientFromConfig()
// 	if err != nil { return err }
// 	flKeys, err := cmd.Flags().GetBool("keys")
// 	if err != nil { return err }
// 	if flKeys {
// 		prefix := arg
// 		return getKeys(cli, prefix)
// 	} else {
// 		key := arg
// 		val, err := cli.Get(key)
// 		if err != nil { return err }
// 		if val == nil { return nil }
// 		fmt.Println(string(val))
// 	}
// 	return nil
// }


// func getKeys (cli *dylt.EtcdClient, prefix string) error {
// 	kvs, err := cli.GetKeys(prefix)
// 	if err != nil { return err }
// 	for _, kv := range kvs {
// 		fmt.Println(kv)
// 	}
// 	return nil
// }