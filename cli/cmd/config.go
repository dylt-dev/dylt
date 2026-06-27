package cmd

import (
	"github.com/dylt-dev/dylt/api"
)

type ConfigOpts struct {
}

func NewConfigCommand(cmdline Cmdline, parent Command) Command {
	// config command
	cfg := BaseCommandConfig[ConfigOpts]{
		name:            "config",
		opts:            ConfigOpts{},
		isUsageOnNoArgs: true,
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
	Key string `pos:"0" impl:"api.RunConfigGet"`
}

func NewConfigGetCommand(cmdline Cmdline, parent Command) *BaseCommand[ConfigGetOpts] {
	// create config object + BaseCommand
	cfg := BaseCommandConfig[ConfigGetOpts]{
		name:      "config.get",
		fnRun:     func(cmd *BaseCommand[ConfigGetOpts]) error { return api.RunConfigGet(cmd.opts.Key) },
		opts:      ConfigGetOpts{},
		usage:     CreateUsageString(USG_Config_Get),
		validator: ArgCountValidator{nExpected: 1},
	}
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any

	// subcommand map if any

	// done
	return cmd
}

type ConfigSetOpts struct {
	Key   string `pos:"0" impl:"api.RunConfigSet"`
	Value string `pos:"1"`
}

func NewConfigSetCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	cfg := BaseCommandConfig[ConfigSetOpts]{
		name:      "config.set",
		fnRun:     func(cmd *BaseCommand[ConfigSetOpts]) error {
			return api.RunConfigSet(cmd.opts.Key, cmd.opts.Value)
		},
		opts:      ConfigSetOpts{},
		usage:     CreateUsageString(USG_Config_Set),
		validator: ArgCountValidator{nExpected: 2},
	}
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any

	// subcommand map if any

	// done
	return cmd
}

type ConfigShowOpts struct {
}

func NewConfigShowCommand(cmdline Cmdline, parent Command) Command {
	// create config object + BaseCommand
	cfg := BaseCommandConfig[ConfigShowOpts]{
		name:      "config.show",
		fnRun:     func(cmd *BaseCommand[ConfigShowOpts]) error { return api.RunConfigShow() },
		opts:      ConfigShowOpts{},
		usage:     CreateUsageString(USG_Config_Get),
		validator: ArgCountValidator{nExpected: 0},
	}
	cmd := NewBaseCommand(cmdline, parent, cfg)

	// flags + args if any

	// subcommand map if any

	// done
	return cmd
}
