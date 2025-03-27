package cmd

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/common"
)

type InitCommand struct {
	*flag.FlagSet
	EtcdDomain string
}

func NewInitCommand () *InitCommand {
	// create command
	flagSet := flag.NewFlagSet("init", flag.ExitOnError)
	cmd := InitCommand { FlagSet: flagSet }
	// init flag vars
	flagSet.StringVar(&cmd.EtcdDomain, "etcd-domain", "", "etcd-domain")
	
	return &cmd
}

func (cmd *InitCommand) HandleArgs(args []string) error {
	// parse flags
	err := cmd.Parse(args)
	if err != nil { return err }
	var requiredFlag string = "etcd-domain"
	if cmd.Lookup(requiredFlag).Value.String() == "" { return fmt.Errorf("required flag missing: %s", requiredFlag)}
	// validate arg count
	cmdArgs := cmd.Args()
	cmdName := "init"
	nExpected := 0
	if len(cmdArgs) != nExpected { return fmt.Errorf("`%s` expects %d argument(s); received %d", cmdName, nExpected, len(cmdArgs)) }
	// init positional params (nop - no params)

	return nil
}

func (cmd *InitCommand) Run (args []string) error {
	slog.Debug("InitCommand.Run()", "args", args)
	// parse flags & get positional args
	err := cmd.HandleArgs(args)
	if err != nil { return err }
	// execute command
	err = RunInit(cmd.EtcdDomain)
	if err != nil { return err }

	return nil
}

func RunInit (etcdDomain string) error {
	slog.Debug("RunInit()", "etcDomain", etcdDomain)
	// create a new config file using the etcdDomain
	if etcdDomain == "" { return errors.New("etcd-domain must be set") }
	cfg := common.ConfigStruct{ EtcdDomain: etcdDomain}
	err := common.SaveConfig(cfg)
	if err != nil { return err }

	return nil
}


// func CreateInitCommand() *cobra.Command {
// 	command := cobra.Command{
// 		Use:   "init",
// 		Short: "Initialize dylt",
// 		Long:  "Initialize dylt",
// 		RunE:  runInitCommand,
// 	}
// 	command.Flags().String("etcd-domain", "", "etcd cluster to activate")
// 	command.MarkFlagRequired("etcd-domain")
// 	return &command
// }

// func runInitCommand(cmd *cobra.Command, args []string) error {
// 	etcdDomain, err := cmd.Flags().GetString("etcd-domain")
// 	if err != nil {
// 		return err
// 	}
// 	initInfo := dylt.InitStruct{
// 		EtcdDomain: etcdDomain,
// 	}
// 	err = dylt.Init(&initInfo)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
