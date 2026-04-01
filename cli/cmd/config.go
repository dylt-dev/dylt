package cmd

import (
	"github.com/dylt-dev/dylt/api"
)

type ConfigCommand struct {
	*BaseCommand
}

func NewConfigCommand(cmdline Cmdline, parent Command) *ConfigCommand {
	// config command
	name := "config"
	cmdMap := CommandMap{
		"get":  ConfigGetCommandF.New,
		"set":  ConfigSetCommandF.New,
		"show": ConfigShowCommandF.New,
	}
	validator := ArgCountGEValidator{nExpected: 0}
	cmd := ConfigCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Config, cmdMap, validator)}
	cmd.isUsageOnNoArgs = true

	// init flag vars (nop -- no flags)

	return &cmd
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
type ConfigGetCommand struct {
	*BaseCommand
	Key string  // arg 0
}

func NewConfigGetCommand(cmdline Cmdline, parent Command) *ConfigGetCommand {
	// config get command
	name := "config.get"
	validator := ArgCountValidator{nExpected: 1}
	cmd := &ConfigGetCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Config_Get, nil, validator)}
	cmd.argmap  = map[int]*string {
		0: &cmd.Key,
	}
	cmd.fnRun = func () error { return api.RunConfigGet(cmd.Key) }
		
	// init flag vars (nop -- no flags)

	return cmd
}

type ConfigSetCommand struct {
	*BaseCommand
	Key   string // arg 0
	Value string // arg 1
}

func NewConfigSetCommand(cmdline Cmdline, parent Command) *ConfigSetCommand {
	// config set command
	name := "config.set"
	validator := ArgCountValidator{nExpected: 2}
	cmd := &ConfigSetCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Config_Set, nil, validator)}
	cmd.argmap  = map[int]*string {
		0: &cmd.Key,
		1: &cmd.Value,
	}
	cmd.fnRun = func () error { return api.RunConfigSet(cmd.Key, cmd.Value) }

	return cmd
}

type ConfigShowCommand struct {
	*BaseCommand
}

func NewConfigShowCommand(cmdline Cmdline, parent Command) *ConfigShowCommand {
	// config show command
	name := "config.set"
	validator := ArgCountValidator{nExpected: 0}
	cmd := &ConfigShowCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Config_Show, nil, validator)}
	cmd.fnRun = func () error { return api.RunConfigShow() }
	//init flags (if any)

	return cmd
}
