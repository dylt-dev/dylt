package cmd

import (
	"log/slog"

	"github.com/dylt-dev/dylt/api"
)

type HostOpts struct {
}

func NewHostCommand(cmdline Cmdline, parent Command) Command {
	// host command
	name := "host"
	opts := HostOpts{}
	cfg := BaseCommandConfig[HostOpts]{
		name:            name,
		isUsageOnNoArgs: true,
		opts:            opts,
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
	Gid int // --gid
	Uid int // --uid
}

func NewHostInitCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	name := "host.init"
	opts := HostInitOpts{}
	fnRun := func(cmd *BaseCommand[HostInitOpts]) error { return api.RunHostInit(cmd.opts.Uid, cmd.opts.Gid) }
	cfg := BaseCommandConfig[HostInitOpts]{
		name:            name,
		fnRun:           fnRun,
		isUsageOnNoArgs: true,
		opts:            opts,
		usage:           CreateUsageString(USG_Config_Get),
		validator:       ArgCountValidator{nExpected: 0},
	}
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any
	cmd.IntVar(&cmd.opts.Gid, "gid", 2000, "gid")
	cmd.IntVar(&cmd.opts.Uid, "uid", 2000, "uid")

	// subcommand map if any

	// done
	return cmd
}
