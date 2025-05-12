package cmd

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/eco"
)

type GetCommand struct {
	*flag.FlagSet
	Key string				// arg 0
}

func NewGetCommand () *GetCommand {
	// create command
	flagSet := flag.NewFlagSet("get", flag.PanicOnError)
	cmd := GetCommand { FlagSet: flagSet }
	// init flag vars (nop - command has no flags)
	
	return &cmd
}

func (cmd *GetCommand) HandleArgs (args []string) error {
	// parse flags
	err := cmd.Parse(args)
	if err != nil { return err }
	// validate arg count
	cmdArgs := cmd.Args()
	cmdName := "get"
	nExpected := 1
	if len(cmdArgs) != nExpected { return fmt.Errorf("`%s` expects %d argument(s); received %d", cmdName, nExpected, len(cmdArgs)) }
	// init positional params
	cmd.Key = cmdArgs[0]

	return nil
}

func (cmd *GetCommand) Run (args []string) error {
	slog.Debug("GetCommand.Run()", "args", args)
	// parse flags & get positional args
	err := cmd.HandleArgs(args)
	if err != nil { return err }
	// execute command
	err = RunGet(cmd.Key)
	if err != nil { return err }

	return nil
}

func RunGet (key string) error {
	slog.Debug("RunGet()", "key", key)
	// create etcd client, get value for key, + output value
	cli, err := eco.CreateEtcdClientFromConfig()
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