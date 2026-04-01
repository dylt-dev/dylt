package cmd

import (
	"log/slog"

	"github.com/dylt-dev/dylt/api"
	"github.com/dylt-dev/dylt/common"
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
	// init flag vars (nop -- no flags)

	return &cmd
}

func (cmd *ConfigCommand) Run() error {
	slog.Debug("ConfigCommand.Run()", "args", cmd.Cmdline)

	// Check for 0 args; if so print usage & return
	if len(cmd.Cmdline) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// Parse flags & get positional args
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

	// Execute command
	err = RunConfig(args, cmd)
	return err
}

func RunConfig(cmdline Cmdline, parent Command) error {
	slog.Debug("RunConfig()", "cmdline", cmdline, "parent", parent)
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
	// init flag vars (nop -- no flags)

	return cmd
}

func (cmd *ConfigGetCommand) Run() error {
	slog.Debug("ConfigGetCommand.Run()", "args", cmd.Cmdline)

	// Parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}

	// If help flag set, print usage + return
	if cmd.Help {
		cmd.PrintUsage()
		return nil
	}

	// Execute command
	err = api.RunConfigGet(cmd.Key)
	return err
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

	return cmd
}

func (cmd *ConfigSetCommand) Run() error {
	slog.Debug("ConfigSetCommand.Run()", "args", cmd.Cmdline)

	// Parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}

	// If help flag set, print usage + return
	if cmd.Help {
		cmd.PrintUsage()
		return nil
	}

	// Execute command
	err = api.RunConfigSet(cmd.Key, cmd.Value)
	return err
}

type ConfigShowCommand struct {
	*BaseCommand
}

func NewConfigShowCommand(cmdline Cmdline, parent Command) *ConfigShowCommand {
	// config show command
	name := "config.set"
	validator := ArgCountValidator{nExpected: 0}
	cmd := &ConfigShowCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Config_Show, nil, validator)}

	//init flags (if any)

	return cmd
}

func (cmd *ConfigShowCommand) Run() error {
	slog.Debug("ConfigShowCommand.Run()", "args", cmd.Cmdline)

	// Parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}

	// If help flag set, print usage + return
	if cmd.Help {
		cmd.PrintUsage()
		return nil
	}

	// Execute command
	err = api.RunConfigShow()
	return err
}
