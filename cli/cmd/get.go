package cmd

import (
	"github.com/dylt-dev/dylt/lib"
)

type GetOpts struct {
	Key string // arg 0
}

func NewGetCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	name := "get"
	opts := GetOpts{}
	fnRun := func (cmd *BaseCommand[GetOpts]) error { return lib.RunGet(cmd.opts.Key) }

	cfg := BaseCommandConfig[GetOpts]{
		name:            name,
		fnRun:           fnRun,
		opts:            opts,
		usage:           USG_Call,
		validator:       ArgCountValidator{nExpected: 1},
	}
	cmd := NewBaseCommand(cmdline, parent, cfg)
	
	// flags + args if any 
	cmd.argMap = ArgMap{
		0: &opts.Key,
	}
	
	// subcommand map if any
	
	// done
	return cmd
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
