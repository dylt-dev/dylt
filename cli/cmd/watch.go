package cmd

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/dylt-dev/dylt/lib"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type WatchCommand struct {
	*flag.FlagSet
	SubCommand string
	SubArgs    []string
}

func NewWatchCommand() *WatchCommand {
	// create command
	flagSet := flag.NewFlagSet("config", flag.ExitOnError)
	cmd := WatchCommand{FlagSet: flagSet}
	// init flag vars (nop -- no flags)

	return &cmd
}

func (cmd *WatchCommand) HandleArgs(args []string) error {
	// parse flags
	err := cmd.Parse(args)
	if err != nil { return err }
	// validate arg count
	cmdArgs := cmd.Args()
	cmdName := "config"
	nExpected := 1
	if len(cmdArgs) < nExpected { return fmt.Errorf("`%s` expects >=%d argument(s); received %d", cmdName, nExpected, len(cmdArgs)) }
	// init positional params
	cmd.SubCommand = cmdArgs[0]
	cmd.SubArgs = cmdArgs[1:]

	return nil
}

func (cmd *WatchCommand) Run(args []string) error {
	slog.Debug("WatchCommand.Run()", "args", args)
	// Parse flags & get positional args
	err := cmd.HandleArgs(args)
	if err != nil { return err }
	// Execute command
	err = RunWatch(cmd.SubCommand, cmd.SubArgs)
	if err != nil { return err }

	return nil
}

func RunWatch(subCommand string, subCmdArgs []string) error {
	slog.Debug("RunWatch()", "subCommand", subCommand, "subCmdArgs", subCmdArgs)
	// Create the subcommand and run it
	subCmd, err := createWatchSubCommand(subCommand)
	if err != nil { return err }
	err = subCmd.Run(subCmdArgs)
	if err != nil { return err }

	return nil
}

func createWatchSubCommand(cmdName string) (Command, error) {
	switch cmdName {
	case "script":
		return NewWatchScriptCommand(), nil
	default:
		return nil, fmt.Errorf("unrecognized command: %s", cmdName)
	}
}

type WatchScriptCommand struct {
	*flag.FlagSet
	ScriptKey string			// arg 0
	TargetPath string			// arg 1
}

func NewWatchScriptCommand() *WatchScriptCommand {
	flagSet := flag.NewFlagSet("config.get", flag.ExitOnError)
	return &WatchScriptCommand{FlagSet: flagSet}
}

func (cmd *WatchScriptCommand) HandleArgs(args []string) error {
	// parse flags
	err := cmd.Parse(args)
	if err != nil { return err }
	// validate arg count
	cmdArgs := cmd.Args()
	cmdName := "watch script"
	nExpected := 2
	if len(cmdArgs) != nExpected { return fmt.Errorf("`%s` expects %d argument(s); received %d", cmdName, nExpected, len(cmdArgs)) }
	// init positional params
	cmd.ScriptKey = cmdArgs[0]
	cmd.TargetPath = cmdArgs[1]

	return nil
}

func (cmd *WatchScriptCommand) Run(args []string) error {
	slog.Debug("WatchScriptCommand.Run()", "args", args)
	// Parse flags & get positional args
	err := cmd.HandleArgs(args)
	if err != nil { return err }
	// Execute command
	err = RunWatchScript(cmd.ScriptKey, cmd.TargetPath)
	if err != nil { return err }

	return nil

}

func RunWatchScript(scriptKey string, targetPath string) error {
	slog.Debug("RunWatchScript()", "scriptKey", scriptKey, targetPath, "targetPath")
	// Get etcd client
	slog.Info("Creating etcd client ...")
	cli, err := lib.CreateEtcdClientFromConfig()
	if err != nil { return err }
	
	// Create watch
	ctx := clientv3.WithRequireLeader(context.Background())
	chWatch := cli.Watch(ctx, scriptKey, clientv3.WithKeysOnly()) 
	
	// Loop over watch channel
	var resp clientv3.WatchResponse
	for resp = range chWatch {
		for _, ev := range resp.Events {
			key := string(ev.Kv.Key)
			fmt.Printf("%s updated\n", key)
			slog.Debug("Update detected", "key", key)
			if key != scriptKey { return fmt.Errorf("key mismatch: execpted %s, got %s", scriptKey, key)}
			ctx := context.Background()
			resp, err := cli.Client.Get(ctx, key)
			if err != nil { return err }
			if len(resp.Kvs) == 0 {
				slog.Warn("No KVs found for watch")
				return nil
			}
			val := (*resp).Kvs[0].Value
			slog.Debug("Value found", "len(val)", len(val))
			fmt.Printf("Writing value to %s ...\n", targetPath)
			err = os.WriteFile(targetPath, val, 0755)
			if err != nil { return err }	
		}
	}

	fmt.Println("Done.")
	return nil
}