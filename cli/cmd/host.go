package cmd

import (
	"log/slog"

	"github.com/dylt-dev/dylt/api"
)

type HostCommand struct {
	*BaseCommand
}

func NewHostCommand(cmdline Cmdline, parent Command) *HostCommand {
	// host command
	name := "host"
	cmdMap := CommandMap{
		"init": HostInitCommandF.New,
	}
	validator := ArgCountGEValidator{nExpected: 0}
	cmd := &HostCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Host, cmdMap, validator)}
	cmd.isUsageOnNoArgs = true
	
	//init flags (if any)
	
	return cmd
}

func RunHost(cmdline Cmdline, parent Command) error {
	slog.Debug("RunHost()", "cmdline", cmdline, "parent", parent)
	// Create the subcommand and run it
	subCmd, err := parent.CreateSubCommand()
	if err != nil {
		return err
	}
	err = subCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

type HostInitCommand struct {
	*BaseCommand
	Gid int // --gid
	Uid int // --uid
}

func NewHostInitCommand(cmdline Cmdline, parent Command) *HostInitCommand {
	// host init command
	name := "host.init"
	validator := ArgCountValidator{nExpected: 0}
	cmd := &HostInitCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Host_Init, nil, validator)}
	cmd.fnRun = func () error { return api.RunHostInit(cmd.Uid, cmd.Gid) }
	
	//init flags (if any)
	cmd.IntVar(&cmd.Gid, "gid", 2000, "gid")
	cmd.IntVar(&cmd.Uid, "uid", 2000, "uid")
	cmd.isUsageOnNoArgs = true

	return cmd
}
