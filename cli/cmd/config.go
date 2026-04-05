package cmd

import (
	"github.com/dylt-dev/dylt/api"
)

type ConfigOpts struct {
}

func NewConfigCommand(cmdline Cmdline, parent Command[ConfigOpts]) Command[ConfigOpts] {
	// config command
	name := "config"
	opts := ConfigOpts{}
	cfg := BaseCommandConfig[ConfigOpts]{
		name:            name,
		isUsageOnNoArgs: true,
		opts:            opts,
		usage:           CreateUsageString(USG_Config),
		validator:       ArgCountGEValidator{nExpected: 0},
	}
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any

	// subcommand map if any
	cmd.subCommandMap = CommandMap{
		"get":  ConfigGetCommandF.New,
		"set":  ConfigSetCommandF.New,
		"show": ConfigShowCommandF.New,
	}

	// done
	return cmd
}

// func RunConfig(cmdline Cmdline, parent Command) error {
// 	slog.Debug("RunConfig()", "cmdline", cmdline, "parent", parent)
// 	// Create the subcommand and run it
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

// Usage
//
//	dylt get key     # get key from config
type ConfigGetOpts struct {
	Key string // arg 0
}

func NewConfigGetCommand(cmdline Cmdline, parent Command) *BaseCommand[ConfigGetOpts] {
	// create config object + BaseCommand
	name := "config.get"
	opts := ConfigGetOpts{}
	fnRun := func(cmd *BaseCommand[ConfigGetOpts]) error { return api.RunConfigGet(cmd.opts.Key) }
	cfg := BaseCommandConfig[ConfigGetOpts]{
		name:      name,
		fnRun:     fnRun,
		opts:      opts,
		usage:     CreateUsageString(USG_Config_Get),
		validator: ArgCountValidator{nExpected: 1},
	}
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any
	cmd.argMap = map[int]*string{
		0: &opts.Key,
	}

	// subcommand map if any

	// done
	return cmd
}

type ConfigSetOpts struct {
	Key   string // arg 0
	Value string // arg 1
}

func NewConfigSetCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	name := "config.set"
	opts := ConfigSetOpts{}
	fnRun := func(cmd *BaseCommand[ConfigSetOpts]) error {
		return api.RunConfigSet(cmd.opts.Key, cmd.opts.Value)
	}
	cfg := BaseCommandConfig[ConfigSetOpts]{
		name:      name,
		fnRun:     fnRun,
		opts:      opts,
		usage:     CreateUsageString(USG_Config_Set),
		validator: ArgCountValidator{nExpected: 2},
	}
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any
	cmd.argMap = map[int]*string{
		0: &opts.Key,
		1: &opts.Value,
	}

	// subcommand map if any

	// done
	return cmd
}

type ConfigShowOpts struct {
}

func NewConfigShowCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	name := "config.show"
	opts := ConfigShowOpts{}
	fnRun := func(cmd *BaseCommand[ConfigShowOpts]) error { return api.RunConfigShow() }
	cfg := BaseCommandConfig[ConfigShowOpts]{
		name:      name,
		fnRun:     fnRun,
		opts:      opts,
		usage:     CreateUsageString(USG_Config_Get),
		validator: ArgCountValidator{nExpected: 0},
	}
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any

	// subcommand map if any

	// done
	return cmd
}
