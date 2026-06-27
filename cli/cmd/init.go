package cmd

import (
	"github.com/dylt-dev/dylt/api"
)

type InitOpts struct {
	EtcdDomain string `flag:"etcd-domain" default:"" desc:"etcd cluster to activate" impl:"api.RunInit"`
}

func NewInitCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	cfg := BaseCommandConfig[InitOpts]{
		name: "init",
		fnRun: func (cmd *BaseCommand[InitOpts]) error { return api.RunInit(cmd.opts.EtcdDomain) },
		opts: InitOpts{},
		usage: CreateUsageString(USG_Init),
		validator: ArgCountValidator{nExpected: 0},
	}	
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any

	// subcommand map if any
	
	// done
	return cmd
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
