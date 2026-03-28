package cmd

import (
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/api"
)

type InitCommand struct {
	*BaseCommand
	EtcdDomain string
}

func NewInitCommand(cmdline Cmdline, parent Command) *InitCommand {
	// init command
	name := "init"
	cmd := &InitCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Init)}
	
	//init flags (if any)
	cmd.FlagSet.StringVar(&cmd.EtcdDomain, "etcd-domain", "", "etcd-domain")

	return cmd
}

func (cmd *InitCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}

	// if Help flag is set, no further processing is necessary
	if cmd.Help {
		return nil
	}

	// validate required flags
	var requiredFlag string = "etcd-domain"
	if cmd.Lookup(requiredFlag).Value.String() == "" {
		cmd.PrintUsage()
		return fmt.Errorf("required flag missing: %s", requiredFlag)
	}
	// validate arg count
	cmdArgs, _ := cmd.Args()
	nExpected := 0
	if len(cmdArgs) != nExpected {
		cmd.PrintUsage()
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
			cmdString,
			nExpected,
			len(cmdArgs))
	}
	// init positional params (nop - no params)

	return nil
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
