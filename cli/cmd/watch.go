package cmd

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/dylt-dev/dylt/common"
	"github.com/dylt-dev/dylt/eco"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type WatchCommand struct {
	*BaseCommand
}

func NewWatchCommand(cmdline Cmdline) *WatchCommand {
	// create command
	flagSet := flag.NewFlagSet("watch", flag.ExitOnError)
	cmd := WatchCommand{BaseCommand: &BaseCommand{Cmdline: cmdline, FlagSet: flagSet}}
	// init flag vars (nop -- no flags)

	return &cmd
}

func (cmd *WatchCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}
	// validate arg count
	var cmdArgs Cmdline = cmd.Args()
	nExpected := 1
	if len(cmdArgs) < nExpected {
		return fmt.Errorf("`%s` expects >=%d argument(s); received %d",
		                  cmd.Cmdline.Command(),
						  nExpected,
						  len(cmdArgs))
	}
	// init positional params (nop - no params)

	return nil
}

func (cmd *WatchCommand) PrintUsage() {
	PrintUsage(USG_Watch)
}

func (cmd *WatchCommand) Run() error {
	slog.Debug("WatchCommand.Run()", "args", cmd.Cmdline)
	// Check for 0 args; if so print usage & return
	if len(cmd.Cmdline) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}
	// Parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}
	// Execute command
	err = RunWatch(cmd.SubCommand(), cmd.SubArgs())
	if err != nil {
		return err
	}

	return nil
}

func RunWatch(subCommand string, subCmdArgs []string) error {
	slog.Debug("RunWatch()", "subCommand", subCommand, "subCmdArgs", subCmdArgs)
	// Create the subcommand and run it
	subCmd, err := createWatchSubCommand(subCommand, subCmdArgs)
	if err != nil {
		return err
	}
	err = subCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func createWatchSubCommand(cmdName string, subCmdArgs Cmdline) (Command, error) {
	switch cmdName {
	case "script":
		return NewWatchScriptCommand(subCmdArgs), nil
	case "svc":
		return NewWatchSvcCommand(subCmdArgs), nil
	default:
		return nil, fmt.Errorf("unrecognized command: %s", cmdName)
	}
}

// Usage
//
//	watch script scriptKey targetPath
type WatchScriptCommand struct {
	*BaseCommand
	ScriptKey  string // arg 0
	TargetPath string // arg 1
}

func NewWatchScriptCommand(cmdline Cmdline) *WatchScriptCommand {
	// create command
	flagSet := flag.NewFlagSet("config.get", flag.ExitOnError)
	cmd := WatchScriptCommand{BaseCommand: &BaseCommand{Cmdline: cmdline, FlagSet: flagSet}}
	// init flag vars (nop -- no flags)

	return &cmd
}

func (cmd *WatchScriptCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}
	// validate arg count
	parentName := "watch"
	cmdName := cmd.Cmdline.Command()
	nExpected := 2
	if len(cmd.Cmdline) != nExpected {
		cmd.PrintUsage()
		return fmt.Errorf("`%s %s` expects %d argument(s); received %d",
		                  parentName,
		                  cmdName,
						  nExpected,
						  len(cmd.Cmdline))
	}
	// init positional params
	cmd.ScriptKey = cmd.Cmdline[0]
	cmd.TargetPath = cmd.Cmdline[1]

	return nil
}

func (cmd *WatchScriptCommand) PrintUsage() {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Watch_Script)
	fmt.Println()
}

func (cmd *WatchScriptCommand) Run() error {
	slog.Debug("WatchScriptCommand.Run()", "args", cmd.Cmdline)
	// Parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}
	// Execute command
	err = RunWatchScript(cmd.ScriptKey, cmd.TargetPath)
	if err != nil {
		return err
	}

	return nil
}

// Usage
//
//	watch svc name
type WatchSvcCommand struct {
	*BaseCommand
	Name string
}

func NewWatchSvcCommand(cmdline Cmdline) *WatchSvcCommand {
	// create command
	flagSet := flag.NewFlagSet("config", flag.ExitOnError)
	cmd := WatchSvcCommand{BaseCommand: &BaseCommand{Cmdline: cmdline, FlagSet: flagSet}}
	// init flag vars (nop -- no flags)

	return &cmd
}

func (cmd *WatchSvcCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}
	// validate arg count
	cmdName := cmd.Cmdline[0]
	nExpected := 1
	if len(cmd.Cmdline) != nExpected {
		cmd.PrintUsage()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
		                  cmdName,
						  nExpected,
						  len(cmd.Cmdline))
	}
	// init positional params
	cmd.Name = cmd.Args()[0]

	return nil
}

func (cmd *WatchSvcCommand) PrintUsage() {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Watch_Svc)
	fmt.Println()
}

func (cmd *WatchSvcCommand) Run() error {
	slog.Debug("WatchSvcCommand.Run()", "args", cmd.Cmdline)
	// Parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}
	// Execute command
	err = RunWatchSvc()
	if err != nil {
		return err
	}

	return nil

}

func RunWatchScript(scriptKey string, targetPath string) error {
	slog.Debug("RunWatchScript()", "scriptKey", scriptKey, targetPath, "targetPath")
	// Get etcd client
	slog.Info("Creating etcd client ...")
	cli, err := eco.CreateEtcdClientFromConfig()
	if err != nil {
		return err
	}

	// Create watch
	ctx := clientv3.WithRequireLeader(context.Background())
	fmt.Printf("Watching %s ...", scriptKey)
	chWatch := cli.Watch(ctx, scriptKey, clientv3.WithKeysOnly())

	// Loop over watch channel
	var resp clientv3.WatchResponse
	for resp = range chWatch {
		for _, ev := range resp.Events {
			key := string(ev.Kv.Key)
			fmt.Printf("%s updated\n", key)
			slog.Debug("Update detected", "key", key)

			if key != scriptKey {
				return fmt.Errorf("key mismatch: execpted %s, got %s", scriptKey, key)
			}

			ctx := context.Background()
			resp, err := cli.Client.Get(ctx, key)
			if err != nil {
				return err
			}
			if len(resp.Kvs) == 0 {
				slog.Warn("No KVs found for watch")
				return nil
			}

			val := (*resp).Kvs[0].Value
			slog.Debug("Value found", "len(val)", len(val))
			fmt.Printf("Writing value to %s ...\n", targetPath)
			err = os.WriteFile(targetPath, val, 0755)
			if err != nil {
				return err
			}

			fmt.Printf("Watching %s ...\n", scriptKey)
		}
	}

	fmt.Println("Done.")
	return nil
}

func RunWatchSvc() error {
	slog.Debug("RunWatchSvc()")

	slog.Info("Creating etcd client ...")
	cli, err := eco.CreateEtcdClientFromConfig()
	if err != nil {
		return err
	}

	// Create watch
	prefix := "svc/"
	ctx := clientv3.WithRequireLeader(context.Background())
	fmt.Printf("Watching %s ...", prefix)
	chWatch := cli.Watch(ctx, prefix, clientv3.WithKeysOnly(), clientv3.WithPrefix())

	// Loop over watch channel
	var resp clientv3.WatchResponse
	for resp = range chWatch {
		for _, ev := range resp.Events {
			key := string(ev.Kv.Key)
			fmt.Printf("%s updated\n", key)
			slog.Debug("Update detected", "key", key)

			fmt.Printf("Watching %s ...\n", prefix)
		}
	}

	fmt.Println("Done.")
	return nil
}
