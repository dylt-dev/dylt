package cmd

import (
	"github.com/dylt-dev/dylt/api"
	"github.com/dylt-dev/dylt/common"
)

type MiscCommand struct {
	*BaseCommand
}

func NewMiscCommand(cmdline Cmdline, parent Command) *MiscCommand {
	// misc command
	name := "misc"
	cmdMap := CommandMap{
		"create-two-node-cluster": CreateTwoNodeClusterCommandF.New,
		"gen-etcd-run-script": GenEtcdRunScriptCommandF.New,
		"lookup": LookupCommandF.New,
	}
	validator := ArgCountGEValidator{nExpected: 0}
	cmd := &MiscCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Misc, cmdMap, validator)}
	
	//init flags (if any)
	
	return cmd
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

	// if CommandMap exists run subcommand
	cmdMap := cmd.CommandMap()
	if cmdMap != nil {
		subCmd, err := cmd.CreateSubCommand()
		if err != nil {
			return err
		}
		err = subCmd.Run()
		return err
	}

	// execute command
	if cmd.fnRun != nil {
		return cmd.fnRun()
	}

	return nil
}

// func RunMisc(cmdline Cmdline, parent Command) error {
// 	common.Logger.Debug("RunMisc()", "cmdline", cmdline, "parent", parent)
// 	// create the subcommand and run it
// 	subCmd, err := parent.CreateSubCommand()
// 	if err != nil {
// 		return err
// 	}
// 	err = subCmd.Run()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }


type CreateTwoNodeClusterCommand struct {
	*BaseCommand
}

func NewCreateTwoNodeClusterCommand(cmdline Cmdline, parent Command) *CreateTwoNodeClusterCommand {
	// misc create-two-node-cluster command
	name := "misc.create-two-node-cluster"
	validator := ArgCountGEValidator{nExpected: 0}
	cmd := &CreateTwoNodeClusterCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Misc_TwoNode, nil, validator)}
	cmd.fnRun = func () error { return api.RunCreateTwoNodeCluster() }

	// init flags (if any)
	
	return cmd
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

	// if CommandMap exists run subcommand
	cmdMap := cmd.CommandMap()
	if cmdMap != nil {
		subCmd, err := cmd.CreateSubCommand()
		if err != nil {
			return err
		}
		err = subCmd.Run()
		return err
	}

	// execute command
	if cmd.fnRun != nil {
		return cmd.fnRun()
	}

	return nil
}

type GenEtcdRunScriptCommand struct {
	*BaseCommand
}

func NewGenEtcdRunScriptCommand(cmdline Cmdline, parent Command) *GenEtcdRunScriptCommand {
	// misc gen-etcd-run-script command
	name := "misc.gen-etcd-run-script"
	validator := ArgCountValidator{nExpected: 0}
	cmd := &GenEtcdRunScriptCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Misc_GenScript, nil, validator)}
	cmd.fnRun = func () error { return api.RunGenEtcdRunScript() }
	
	//init flags (if any)
	
	return cmd
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

	// if CommandMap exists run subcommand
	cmdMap := cmd.CommandMap()
	if cmdMap != nil {
		subCmd, err := cmd.CreateSubCommand()
		if err != nil {
			return err
		}
		err = subCmd.Run()
		return err
	}

	// execute command
	if cmd.fnRun != nil {
		return cmd.fnRun()
	}

	return nil
}

type LookupCommand struct {
	*BaseCommand
	Hostname string		//arg 0
}

func NewLookupCommand(cmdline Cmdline, parent Command) *LookupCommand {
	// misc lookup command
	name := "misc.lookup"
	validator := ArgCountValidator{nExpected: 1}
	cmd := &LookupCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Misc_Lookup, nil, validator)}
	cmd.argmap  = map[int]*string {
		0: &cmd.Hostname,
	}
	cmd.fnRun = func () error { return api.RunLookupCommand(cmd.Hostname) }
	
	//init flags (if any)
	
	return cmd
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

	// if CommandMap exists run subcommand
	cmdMap := cmd.CommandMap()
	if cmdMap != nil {
		subCmd, err := cmd.CreateSubCommand()
		if err != nil {
			return err
		}
		err = subCmd.Run()
		return err
	}

	// execute command
	if cmd.fnRun != nil {
		return cmd.fnRun()
	}

	return nil
}
