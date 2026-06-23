package cmd

import (
	"log/slog"

	"github.com/dylt-dev/dylt/api"
)

type HostOpts struct {
}

func NewHostCommand(cmdline Cmdline, parent Command) Command {
	// host command
	cfg := BaseCommandConfig[HostOpts]{
		name:            "host",
		opts:            HostOpts{},
		isUsageOnNoArgs: true,
		usage:           CreateUsageString(USG_Host),
		validator:       ArgCountGEValidator{nExpected: 0},
	}
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any

	// subcommand map if any
	cmd.subCommandMap = CommandMap{
		"init": HostInitCommandF.New,
	}

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

type HostInitOpts struct {
	Gid int `flag:"gid" default:"2000" desc:"group ID" impl:"api.RunHostInit"`
	Uid int `flag:"uid" default:"2000" desc:"user ID"`
}

func NewHostInitCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	cfg := BaseCommandConfig[HostInitOpts]{
		name:            "host.init",
		fnRun:           func(cmd *BaseCommand[HostInitOpts]) error { return api.RunHostInit(cmd.opts.Uid, cmd.opts.Gid) },
		opts:            HostInitOpts{},
		isUsageOnNoArgs: true,
		usage:           CreateUsageString(USG_Config_Get),
		validator:       ArgCountValidator{nExpected: 0},
	}
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any

	// subcommand map if any

	// done
	return cmd
}
