package cmd

import (
	"github.com/dylt-dev/dylt/api"
)

type InitOpts struct {
	EtcdDomain string // --etcd-domain
}

func NewInitCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	name := "init"
	opts := InitOpts{}
	fnRun := func (cmd *BaseCommand[InitOpts]) error { return api.RunInit(cmd.opts.EtcdDomain) }
	cfg := BaseCommandConfig[InitOpts]{
		name: name,
		fnRun: fnRun,
		opts: opts,
		usage: CreateUsageString(USG_Config_Get),
		validator: ArgCountValidator{nExpected: 0},
	}	
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any 
	cmd.FlagSet.StringVar(&cmd.opts.EtcdDomain, "etcd-domain", "", "etcd-domain")

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
