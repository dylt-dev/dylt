package cmd

import (
	"log/slog"

	"github.com/dylt-dev/dylt/api"
)

type InitCommand struct {
	*BaseCommand
	EtcdDomain string // --etcd-domain
}

func NewInitCommand(cmdline Cmdline, parent Command) *InitCommand {
	// init command
	name := "init"
	validator := ArgCountValidator{nExpected: 0}
	cmd := &InitCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Init, nil, validator)}
	
	//init flags (if any)
	cmd.FlagSet.StringVar(&cmd.EtcdDomain, "etcd-domain", "", "etcd-domain")

	return cmd
}

func (cmd *InitCommand) Run() error {
	slog.Debug("InitCommand.Run()", "args", cmd.Cmdline)

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
	err = api.RunInit(cmd.EtcdDomain)
	if err != nil {
		return err
	}

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
