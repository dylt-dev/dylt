package cmd

import (
	"context"
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

func NewWatchCommand(cmdline Cmdline, parent Command) *WatchCommand {
	// watch command
	name := "watch"
	cmd := &WatchCommand{BaseCommand: NewBaseCommand(name, cmdline, parent)}
	
	//init flags (if any)
	
	return cmd
}

func (cmd *WatchCommand) CreateSubCommand() (Command, error) {
	args, flag := cmd.Args()
	if !flag {
		return nil, nil
	}
	return createWatchSubCommand(args, cmd)
}

func (cmd *WatchCommand) HandleArgs() error {
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
	nExpected := 0
	if len(cmdArgs) < nExpected {
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("`%s` expects >=%d argument(s); received %d",
			cmdString,
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

	// If help flag set, print usage + return
	if cmd.Help {
		cmd.PrintUsage()
		return nil
	}

	// If no args, print usage
	args, _ := cmd.Args()
	if len(args) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// Execute command
	err = RunWatch(args, cmd)
	return err
}

func RunWatch(cmdline Cmdline, parent Command) error {
	slog.Debug("RunWatch()", "cmdline", cmdline, "parent", parent)
	// Create the subcommand and run it
	subCmd, err := createWatchSubCommand(cmdline, parent)
	if err != nil {
		return err
	}
	err = subCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func createWatchSubCommand(cmdline Cmdline, parent Command) (Command, error) {
	cmdName := cmdline.Command()
	switch cmdName {
	case "script":
		return NewWatchScriptCommand(cmdline, parent), nil
	case "svc":
		return NewWatchSvcCommand(cmdline, parent), nil
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

func NewWatchScriptCommand(cmdline Cmdline, parent Command) *WatchScriptCommand {
	// watch script command
	name := "watch.script"
	cmd := &WatchScriptCommand{BaseCommand: NewBaseCommand(name, cmdline, parent)}
	
	//init flags (if any)
	
	return cmd
}

func (cmd *WatchScriptCommand) HandleArgs() error {
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
	nExpected := 2
	if len(cmdArgs) != nExpected {
		cmd.PrintUsage()
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
			cmdString,
			nExpected,
			len(cmd.Cmdline))
	}

	// init positional params
	cmd.ScriptKey = cmdArgs[0]
	cmd.TargetPath = cmdArgs[1]

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

	// if Help flag is set, no further processing is necessary
	if cmd.Help {
		return nil
	}

	// Execute command
	err = RunWatchScript(cmd.ScriptKey, cmd.TargetPath)
	return err
}

// Usage
//
//	watch svc name
type WatchSvcCommand struct {
	*BaseCommand
	Name string
}

func NewWatchSvcCommand(cmdline Cmdline, parent Command) *WatchSvcCommand {
	// watch svc command
	name := "watch.svc"
	cmd := &WatchSvcCommand{BaseCommand: NewBaseCommand(name, cmdline, parent)}
	
	//init flags (if any)
	
	return cmd
}

func (cmd *WatchSvcCommand) HandleArgs() error {
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
	cmd.Name = cmdArgs[0]

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

	// If help flag set, print usage + return
	if cmd.Help {
		cmd.PrintUsage()
		return nil
	}

	// If no args, print usage
	if len(cmd.Cmdline) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// Execute command
	err = RunWatchSvc()
	return err

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
