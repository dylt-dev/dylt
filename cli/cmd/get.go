package cmd

import (
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/lib"
)

type GetCommand struct {
	*BaseCommand
	Key string // arg 0
}

func NewGetCommand(cmdline Cmdline, parent Command) *GetCommand {
	// get command
	name := "get"
	cmd := &GetCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Get, nil)}
	
	//init flags (if any)
	
	return cmd
}

func (cmd *GetCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}

	// if Help flag is set, no further processing is necessary
	if cmd.Help {
		return nil
	}

	// validate arg count
	cmdArgs, _ := cmd.Args()
	nExpected := 1
	if len(cmdArgs) != nExpected {
		cmd.PrintUsage()
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
			cmdString,
			nExpected,
			len(cmdArgs))
	}
	// init positional params
	cmd.Key = cmdArgs[0]

	return nil
}

func (cmd *GetCommand) Run() error {
	slog.Debug("GetCommand.Run()", "args", cmd.Cmdline)

	// parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}

	// If help flag set, print usage + return
	if cmd.Help {
		cmd.PrintUsage()
		return nil
	}

	// execute command
	err = lib.RunGet(cmd.Key)
	if err != nil {
		return err
	}

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
