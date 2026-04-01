package cmd

import (
	"log/slog"

	"github.com/dylt-dev/dylt/api"
	"github.com/dylt-dev/dylt/common"
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
	cmd.fnRun = func () error { return api.RunInit(cmd.EtcdDomain) }
	
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

	// Check for 0 args; if so print usage & return
	args, _ := cmd.Args()
	if len(args) == 0 && cmd.UsageOnNoArgs() {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// If CommandMap exists run subcommand
	cmdMap := cmd.CommandMap()
	if cmdMap != nil {
		subCmd, err := cmd.CreateSubCommand()
		if err != nil {
			return err
		}
		err = subCmd.Run()
		return err
	}

	// Execute command
	if cmd.fnRun != nil {
		return cmd.fnRun()
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
