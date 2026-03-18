package cmd

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/dylt-dev/dylt/common"
)

type ConfigCommand struct {
	*BaseCommand
	SubCommand string
	SubArgs    []string
}

func NewConfigCommand(cmdline Cmdline) *ConfigCommand {
	// create command
	flagSet := flag.NewFlagSet("config", flag.ExitOnError)
	cmd := ConfigCommand{BaseCommand: &BaseCommand{Cmdline: cmdline, FlagSet: flagSet}}
	// init flag vars (nop -- no flags)

	return &cmd
}

func (cmd *ConfigCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}
	// validate arg count
	cmdArgs, _ := cmd.Args()
	nExpected := 1
	if len(cmdArgs) < nExpected {
		return fmt.Errorf("`%s` expects >=%d argument(s); received %d",
		                  cmd.GetCommandString(),
						  nExpected,
						  len(cmdArgs))
	}

	return nil
}

func (cmd *ConfigCommand) PrintUsage() {
	PrintUsage(USG_Config)
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
	// Execute command
	err = RunConfig(cmd.SubCommand, cmd.SubArgs)
	if err != nil {
		return err
	}

	return nil
}

func RunConfig(subCommand string, subCmdArgs Cmdline) error {
	slog.Debug("RunConfig()", "subCommand", subCommand, "subCmdArgs", subCmdArgs)
	// Create the subcommand and run it
	subCmd, err := createConfigSubCommand(subCommand, subCmdArgs)
	if err != nil {
		return err
	}
	err = subCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func createConfigSubCommand(cmdName string, cmdline Cmdline) (Command, error) {
	switch cmdName {
	case "get":
		return NewConfigGetCommand(cmdline), nil
	case "set":
		return NewConfigSetCommand(cmdline), nil
	case "show":
		return NewConfigShowCommand(cmdline), nil
	default:
		{
			var this *ConfigCommand = nil
			this.PrintUsage()
			return nil, fmt.Errorf("unrecognized command: %s", cmdName)
		}
	}
}

// Usage
//
//	dylt get key     # get key from config
type ConfigGetCommand struct {
	*BaseCommand
	Key string
}

func NewConfigGetCommand(cmdline Cmdline) *ConfigGetCommand {
	flagSet := flag.NewFlagSet("config.get", flag.ExitOnError)
	return &ConfigGetCommand{BaseCommand: &BaseCommand{Cmdline: cmdline, FlagSet: flagSet}}
}

func (cmd *ConfigGetCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}
	// validate arg count
	cmdArgs, _ := cmd.Args()
	if len(cmdArgs) != 1 {
		cmd.PrintUsage()
		return fmt.Errorf("`%s` expects 1 argument(s); received %d",
		                   cmd.GetCommandString(),
		                   len(cmdArgs))
	}
	// init positional params
	cmd.Key = cmdArgs[0]

	return nil
}

func (cmd *ConfigGetCommand) PrintUsage() {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Config_Get)
	fmt.Println()
}

func (cmd *ConfigGetCommand) Run() error {
	slog.Debug("ConfigGetCommand.Run()", "args", cmd.Cmdline)
	// Parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}
	// Execute command
	err = RunConfigGet(cmd.Key)
	if err != nil {
		return err
	}

	return nil

}

func RunConfigGet(key string) error {
	slog.Debug("RunConfigGet()", "key", key)
	val, err := common.GetConfigValue(key)
	if err != nil {
		return err
	}
	fmt.Printf("\n%s\n", val)

	return nil
}

type ConfigSetCommand struct {
	*BaseCommand
	Key   string // arg 0
	Value string // arg 1
}

func NewConfigSetCommand(cmdline Cmdline) *ConfigSetCommand {
	flagSet := flag.NewFlagSet("config.set", flag.ExitOnError)
	return &ConfigSetCommand{BaseCommand: &BaseCommand{Cmdline: cmdline, FlagSet: flagSet}}
}

func (cmd *ConfigSetCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}
	// validate arg count
	cmdArgs, _ := cmd.Args()
	if len(cmdArgs) != 2 {
		cmd.PrintUsage()
		return fmt.Errorf("`%s` expects 2 argument(s); received %d",
		                  cmd.GetCommandString(),
						  len(cmdArgs))
	}
	// init positional params
	cmd.Key = cmdArgs[0]
	cmd.Value = cmdArgs[1]

	return nil
}

func (cmd *ConfigSetCommand) PrintUsage() {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Config_Set)
	fmt.Println()
}

func (cmd *ConfigSetCommand) Run() error {
	slog.Debug("ConfigSetCommand.Run()", "args", cmd.Cmdline)
	// Parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}
	// Execute command
	err = RunConfigSet(cmd.Key, cmd.Value)
	if err != nil {
		return err
	}

	return nil
}

func RunConfigSet(key string, val string) error {
	slog.Debug("RunConfigSet()", "key", key, "val", val)

	// Open the dylt config file for read+write. Create if necessasry.
	cfgFilePath := common.GetConfigFilePath()
	slog.Debug("Opening config file", "cfgFilePath", cfgFilePath)
	f, err := os.OpenFile(cfgFilePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return common.NewError(err)
	}
	defer f.Close()

	// Read the dylt config file as YAML
	data, err := common.ReadYaml(f)
	if err != nil {
		return err
	}

	// Truncate the file to 0 and rewrite
	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	// Set the config map value and write the updated config map
	data.Set(key, val)
	err = common.WriteConfig(data)
	if err != nil {
		return err
	}
	err = common.WriteYaml(data, f)
	if err != nil {
		return err
	}

	return nil
}

type ConfigShowCommand struct {
	*BaseCommand
}

func NewConfigShowCommand(cmdline Cmdline) *ConfigShowCommand {
	flagSet := flag.NewFlagSet("config.show", flag.ExitOnError)
	return &ConfigShowCommand{BaseCommand: &BaseCommand{Cmdline: cmdline, FlagSet: flagSet}}
}

func (cmd *ConfigShowCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}
	// validate arg count
	cmdArgs, _ := cmd.Args()
	nExpected := 0
	if len(cmdArgs) != nExpected {
		cmd.PrintUsage()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
		                  cmd.GetCommandString(),
						  nExpected,
						  len(cmdArgs))
	}
	// init positional params (nop - no positional params)

	return nil
}

func (cmd *ConfigShowCommand) PrintUsage() {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Config_Show)
	fmt.Println()
}

func (cmd *ConfigShowCommand) Run() error {
	slog.Debug("ConfigShowCommand.Run()", "args", cmd.Cmdline)
	// Parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}
	// Execute command
	err = RunConfigShow()
	if err != nil {
		return err
	}

	return nil
}

func RunConfigShow() error {
	err := common.ShowConfig(os.Stdout)
	if err != nil {
		return err
	}

	return nil
}
