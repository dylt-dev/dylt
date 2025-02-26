package cmd

import (
	"errors"
	"flag"

	"github.com/dylt-dev/dylt/lib"
)

type InitCommand struct {
	*flag.FlagSet
	EtcdDomain string
}

func NewInitCommand () *InitCommand {
	flagSet := flag.NewFlagSet("init", flag.PanicOnError)
	cmd := InitCommand { FlagSet: flagSet }
	flagSet.StringVar(&cmd.EtcdDomain, "etcd-domain", "", "etcd-domain")
	return &cmd
}

func (cmd *InitCommand) Run (args []string) error {
	err := cmd.Parse(args)
	if err != nil { return err }
	if cmd.EtcdDomain == "" {
		return errors.New("etcd-domain must be set")
	}
	cfg := lib.ConfigStruct{ EtcdDomain: cmd.EtcdDomain}
	err = lib.SaveConfig(cfg)
	return err
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
