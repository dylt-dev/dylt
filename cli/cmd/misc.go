package cmd

import (
	"github.com/dylt-dev/dylt/api"
)

type MiscCommand struct {
	*BaseCommand
}

func NewMiscCommand(cmdline Cmdline, parent Command) *MiscCommand {
	// misc command
	name := "misc"
	cmdMap := CommandMap{
		"create-two-node-cluster": CreateTwoNodeClusterCommandF.New,
		"gen-etcd-run-script":     GenEtcdRunScriptCommandF.New,
		"lookup":                  LookupCommandF.New,
	}
	validator := ArgCountGEValidator{nExpected: 0}
	cmd := &MiscCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Misc, cmdMap, validator)}
	cmd.isUsageOnNoArgs = true

	//init flags (if any)

	return cmd
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
	cmd.fnRun = func() error { return api.RunCreateTwoNodeCluster() }

	// init flags (if any)

	return cmd
}

type GenEtcdRunScriptCommand struct {
	*BaseCommand
}

func NewGenEtcdRunScriptCommand(cmdline Cmdline, parent Command) *GenEtcdRunScriptCommand {
	// misc gen-etcd-run-script command
	name := "misc.gen-etcd-run-script"
	validator := ArgCountValidator{nExpected: 0}
	cmd := &GenEtcdRunScriptCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Misc_GenScript, nil, validator)}
	cmd.fnRun = func() error { return api.RunGenEtcdRunScript() }

	//init flags (if any)

	return cmd
}

type LookupCommand struct {
	*BaseCommand
	Hostname string //arg 0
}

func NewLookupCommand(cmdline Cmdline, parent Command) *LookupCommand {
	// misc lookup command
	name := "misc.lookup"
	validator := ArgCountValidator{nExpected: 1}
	cmd := &LookupCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Misc_Lookup, nil, validator)}
	cmd.argMap = map[int]*string{
		0: &cmd.Hostname,
	}
	cmd.fnRun = func() error { return api.RunLookupCommand(cmd.Hostname) }

	//init flags (if any)

	return cmd
}
