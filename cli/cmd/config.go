package cmd

import (
	"errors"
	"fmt"
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
	cmd := ConfigCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Config)}
	// init flag vars (nop -- no flags)

	return &cmd
}

func (cmd *ConfigCommand) CreateSubCommand() (Command, error) {
	args, is := cmd.Args()
	if !is {
		return nil, errors.New("Command not Parse()'d")
	}
	return createConfigSubCommand(args, cmd)
}

func (cmd *ConfigCommand) HandleArgs() error {
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
	if len(cmdArgs) < nExpected {
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("`%s` expects >=%d argument(s); received %d",
			cmdString,
			nExpected,
			len(cmdArgs))
	}

	return nil
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
	subCmd, err := createConfigSubCommand(cmdline, parent)
	if err != nil {
		return err
	}
	err = subCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func createConfigSubCommand(cmdline Cmdline, parent Command) (Command, error) {
	cmdName := cmdline.Command()
	cmdMap := CommandMap{
		"get":  ConfigGetCommandF.New,
		"set": ConfigSetCommandF.New,
		"show": ConfigShowCommandF.New,
	}
	
	cmdFactoryFunc, ok := cmdMap[cmdName]
	if !ok {
		parent.PrintUsage()
		return nil, fmt.Errorf("unrecognized command: %s", cmdName)
	}
		
	cmd := cmdFactoryFunc(cmdline, parent)
	return cmd, nil
}

// Usage
//
//	dylt get key     # get key from config
type ConfigGetCommand struct {
	*BaseCommand
	Key string
}

func NewConfigGetCommand(cmdline Cmdline, parent Command) *ConfigGetCommand {
	// config get command
	name := "config.get"
	cmd := &ConfigGetCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Config_Get)}
	// init flag vars (nop -- no flags)

	return cmd
}

func (cmd *ConfigGetCommand) HandleArgs() error {
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
	if len(cmdArgs) != 1 {
		cmd.PrintUsage()
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("`%s` expects 1 argument(s); received %d",
			cmdString,
			len(cmdArgs))
	}

	// init positional params
	cmd.Key = cmdArgs[0]

	return nil
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
	cmd := &ConfigSetCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Config_Set)}

	return cmd
}

func (cmd *ConfigSetCommand) HandleArgs() error {
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
	if len(cmdArgs) != 2 {
		cmd.PrintUsage()
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("`%s` expects 2 argument(s); received %d",
			cmdString,
			len(cmdArgs))
	}

	// init positional params
	cmd.Key = cmdArgs[0]
	cmd.Value = cmdArgs[1]

	return nil
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
	cmd := &ConfigShowCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Config_Show)}

	//init flags (if any)

	return cmd
}

func (cmd *ConfigShowCommand) HandleArgs() error {
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

	// init positional params (nop - no positional params)

	return nil
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
