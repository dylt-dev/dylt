package cmd

import (
	"fmt"

	"github.com/dylt-dev/dylt/api"
	"github.com/dylt-dev/dylt/common"
)

type MiscCommand struct {
	*BaseCommand
}

func NewMiscCommand(cmdline Cmdline, parent Command) *MiscCommand {
	// misc command
	name := "misc"
	cmd := &MiscCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Misc)}
	
	//init flags (if any)
	
	return cmd
}

func (cmd *MiscCommand) CreateSubCommand () (Command, error) {
	args, is := cmd.Args()
	if !is {
		return nil, nil
	}
	return createMiscSubCommand(args, cmd)
}

func (cmd *MiscCommand) HandleArgs() error {
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
	nExpected := 0
	if len(cmd.Cmdline) < nExpected {
		cmd.PrintUsage()
		return fmt.Errorf("`%s` expects >=%d argument(s); received %d",
			cmd.Cmdline[0],
			nExpected,
			len(cmd.Cmdline))
	}

	// init positional params

	return nil
}

func (cmd *MiscCommand) Run() error {
	common.Logger.Debug("MiscCommand.Run()", "args", cmd.Cmdline)

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

	// If no args, print usage
	args, _ := cmd.Args()
	if len(args) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// execute command
	err = RunMisc(args, cmd)
	return err
}

func RunMisc(cmdline Cmdline, parent Command) error {
	common.Logger.Debug("RunMisc()", "cmdline", cmdline, "parent", parent)
	// create the subcommand and run it
	subCmd, err := createMiscSubCommand(cmdline, parent)
	if err != nil {
		return err
	}
	err = subCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func createMiscSubCommand(cmdline Cmdline, parent Command) (Command, error) {
	cmdName := cmdline.Command()
	cmdMap := CommandMap{
		"create-two-node-cluster": CreateTwoNodeClusterCommandF.New,
		"gen-etcd-run-script": GenEtcdRunScriptCommandF.New,
		"lookup": LookupCommandF.New,
	}
	
	cmdFactoryFunc, ok := cmdMap[cmdName]
	if !ok {
		parent.PrintUsage()
		return nil, fmt.Errorf("unrecognized command: %s", cmdName)
	}
		
	cmd := cmdFactoryFunc(cmdline, parent)
	return cmd, nil
	
	// cmdName := cmdline.Command()
	// switch cmdName {
	// case "create-two-node-cluster":
	// 	return CreateTwoNodeClusterCommandF.New(cmdline, parent), nil
	// case "gen-etcd-run-script":
	// 	return GenEtcdRunScriptCommandF.New(cmdline, parent), nil
	// case "lookup":
	// 	return LookupCommandF.New(cmdline, parent), nil
	// default:
	// 	parent.PrintUsage()
	// 	return nil, fmt.Errorf("unrecognized subcommand: %s", cmdName)
	// }
}

type CreateTwoNodeClusterCommand struct {
	*BaseCommand
}

func NewCreateTwoNodeClusterCommand(cmdline Cmdline, parent Command) *CreateTwoNodeClusterCommand {
	// misc create-two-node-cluster command
	name := "misc.create-two-node-cluster"
	cmd := &CreateTwoNodeClusterCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Misc_TwoNode)}
	
	//init flags (if any)
	
	return cmd
}

func (cmd *CreateTwoNodeClusterCommand) HandleArgs() error {
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
	if len(cmd.Cmdline) < nExpected {
		cmd.PrintUsage()
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
			cmdString,
			nExpected,
			len(cmdArgs))
	}

	return nil
}

func (cmd *CreateTwoNodeClusterCommand) Run() error {
	common.Logger.Debug("CreateTwoNodeClusterCommand.Run()", "args", cmd.Cmdline)

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
	// @getit
	err = api.RunCreateTwoNodeCluster()
	return err
}

type GenEtcdRunScriptCommand struct {
	*BaseCommand
}

func NewGenEtcdRunScriptCommand(cmdline Cmdline, parent Command) *GenEtcdRunScriptCommand {
	// misc gen-etcd-run-script command
	name := "misc.gen-etcd-run-script"
	cmd := &GenEtcdRunScriptCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Misc_GenScript)}
	
	//init flags (if any)
	
	return cmd
}

func (cmd *GenEtcdRunScriptCommand) HandleArgs() error {
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
	if len(cmdArgs) != nExpected {
		cmd.PrintUsage()
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
			cmdString,
			nExpected,
			len(cmdArgs))
	}

	return nil
}

func (cmd *GenEtcdRunScriptCommand) Run() error {
	common.Logger.Debug("GenEtcdRunScriptCommand.Run()", "args", cmd.Cmdline)

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
	// @getit
	err = api.RunGenEtcdRunScript()
	if err != nil {
		return err
	}

	common.Logger.WithGroup("g")
	common.Logger.With("bar", "thirteen")
	common.Logger.Debug("testing logger", "foo", "13")
	return nil
}

type LookupCommand struct {
	*BaseCommand
	Hostname string
}

func NewLookupCommand(cmdline Cmdline, parent Command) *LookupCommand {
	// misc lookup command
	name := "misc.lookup"
	cmd := &LookupCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Misc_Lookup)}
	
	//init flags (if any)
	
	return cmd
}

func (cmd *LookupCommand) HandleArgs() error {
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
	cmd.Hostname = cmdArgs[0]

	return nil
}

func (cmd *LookupCommand) Run() error {
	common.Logger.Debug("LookupCommand.Run()", "args", cmd.Cmdline)

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

	// If no args, print usage
	if len(cmd.Cmdline) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// execute command
	// @getit
	err = api.RunLookupCommand(cmd.Hostname)
	return err
}
