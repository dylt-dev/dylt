package cmd

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/dylt-dev/dylt/common"
)

type ConfigCommand struct {
	*flag.FlagSet
	SubCommand string
	SubArgs    []string
}

func NewConfigCommand() *ConfigCommand {
	// create command
	flagSet := flag.NewFlagSet("config", flag.ExitOnError)
	cmd := ConfigCommand{FlagSet: flagSet}
	// init flag vars (nop -- no flags)

	return &cmd
}

func (cmd *ConfigCommand) HandleArgs(args []string) error {
	// parse flags
	err := cmd.Parse(args)
	if err != nil { return err }
	// validate arg count
	cmdArgs := cmd.Args()
	cmdName := "config"
	nExpected := 1
	if len(cmdArgs) < nExpected { return fmt.Errorf("`%s` expects >=%d argument(s); received %d", cmdName, nExpected, len(cmdArgs)) }
	// init positional params
	cmd.SubCommand = cmdArgs[0]
	cmd.SubArgs = cmdArgs[1:]

	return nil
}


func (cmd *ConfigCommand) PrintUsage () {
	PrintMultilineUsage(USG_Config)
}

func (cmd *ConfigCommand) Run(args []string) error {
	slog.Debug("ConfigCommand.Run()", "args", args)
	// Check for 0 args; if so print usage & return
	if len(args) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}
	// Parse flags & get positional args
	err := cmd.HandleArgs(args)
	if err != nil { return err }
	// Execute command
	err = RunConfig(cmd.SubCommand, cmd.SubArgs)
	if err != nil { return err }

	return nil
}

func RunConfig(subCommand string, subCmdArgs []string) error {
	slog.Debug("RunConfig()", "subCommand", subCommand, "subCmdArgs", subCmdArgs)
	// Create the subcommand and run it
	subCmd, err := createConfigSubCommand(subCommand)
	if err != nil { return err }
	err = subCmd.Run(subCmdArgs)
	if err != nil { return err }

	return nil
}

func createConfigSubCommand(cmdName string) (Command, error) {
	switch cmdName {
	case "get": return NewConfigGetCommand(), nil
	case "set": return NewConfigSetCommand(), nil
	case "show": return NewConfigShowCommand(), nil
	default: {
		var this *ConfigCommand = nil	
		this.PrintUsage()
		return nil, fmt.Errorf("unrecognized command: %s", cmdName)
	}
}
}

// Usage
//
//     dylt get key     # get key from config
type ConfigGetCommand struct {
	*flag.FlagSet
	Key string
}

func NewConfigGetCommand() *ConfigGetCommand {
	flagSet := flag.NewFlagSet("config.get", flag.ExitOnError)
	return &ConfigGetCommand{FlagSet: flagSet}
}

func (cmd *ConfigGetCommand) HandleArgs(args []string) error {
	// parse flags
	err := cmd.Parse(args)
	if err != nil { return err }
	// validate arg count
	cmdArgs := cmd.Args()
	if len(cmdArgs) != 1 {
		cmd.PrintUsage()
		return fmt.Errorf("`config get` expects 1 argument(s); received %d", len(cmdArgs)) }
	// init positional params
	cmd.Key = cmdArgs[0]

	return nil
}

func (cmd *ConfigGetCommand) PrintUsage () {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Config_Get)
	fmt.Println()
}

func (cmd *ConfigGetCommand) Run(args []string) error {
	slog.Debug("ConfigGetCommand.Run()", "args", args)
	// Parse flags & get positional args
	err := cmd.HandleArgs(args)
	if err != nil { return err }
	// Execute command
	err = RunConfigGet(cmd.Key)
	if err != nil { return err }

	return nil

}

func RunConfigGet(key string) error {
	slog.Debug("RunConfigGet()", "key", key)
	val, err := common.GetConfigValue(key)	
	if err != nil { return err }
	fmt.Printf("\n%s\n", val)
	
	return nil
}

type ConfigSetCommand struct {
	*flag.FlagSet
	Key string	 		// arg 0
	Value string			// arg 1
}

func NewConfigSetCommand() *ConfigSetCommand {
	flagSet := flag.NewFlagSet("config.set", flag.ExitOnError)
	return &ConfigSetCommand{FlagSet: flagSet}
}

func (cmd *ConfigSetCommand) HandleArgs(args []string) error {
	// parse flags
	err := cmd.Parse(args)
	if err != nil { return err }
	// validate arg count
	cmdArgs := cmd.Args()
	if len(cmdArgs) != 2 {
		cmd.PrintUsage()
		return fmt.Errorf("`config set` expects 2 argument(s); received %d", len(cmdArgs))
	}
	// init positional params
	cmd.Key = cmdArgs[0]
	cmd.Value = cmdArgs[1]

	return nil
}

func (cmd *ConfigSetCommand) PrintUsage () {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Config_Set)
	fmt.Println()
}

func (cmd *ConfigSetCommand) Run(args []string) error {
	slog.Debug("ConfigSetCommand.Run()", "args", args)
	// Parse flags & get positional args
	err := cmd.HandleArgs(args)
	if err != nil { return err }
	// Execute command
	err = RunConfigSet(cmd.Key, cmd.Value)
	if err != nil { return err }

	return nil
}

func RunConfigSet (key string, val string) error {
	slog.Debug("RunConfigSet()", "key", key, "val", val)
	
	// Open the dylt config file for read+write. Create if necessasry.
	cfgFilePath := common.GetConfigFilePath()
	slog.Debug("Opening config file", "cfgFilePath", cfgFilePath)
	f, err := os.OpenFile(cfgFilePath, os.O_CREATE | os.O_RDWR, 0644)
	if err != nil { return common.NewError(err) }
	defer f.Close()

	// Read the dylt config file as YAML
	data, err := common.ReadYaml(f)
	if err != nil { return err }

	// Truncate the file to 0 and rewrite 
	err = f.Truncate(0)
	if err != nil { return err }
	_, err = f.Seek(0, 0)
	if err != nil { return err }
	
	// Set the config map value and write the updated config map
	data.Set(key, val)
	err = common.WriteConfig(data)
	if err != nil { return err }
	err = common.WriteYaml(data, f)
	if err != nil { return err }

	return nil
}

type ConfigShowCommand struct {
	*flag.FlagSet
}

func NewConfigShowCommand() *ConfigShowCommand {
	flagSet := flag.NewFlagSet("config.show", flag.ExitOnError)
	return &ConfigShowCommand{FlagSet: flagSet}
}

func (cmd *ConfigShowCommand) HandleArgs(args []string) error {
	// parse flags
	err := cmd.Parse(args)
	if err != nil { return err }
	// validate arg count
	cmdArgs := cmd.Args()
	cmdName := "config show"
	nExpected := 0
	if len(cmdArgs) != nExpected {
		cmd.PrintUsage()
		return fmt.Errorf("`%s` expects %d argument(s); received %d", cmdName, nExpected, len(cmdArgs))
	}
	// init positional params (nop - no positional params)

	return nil
}

func (cmd *ConfigShowCommand) PrintUsage () {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Config_Show)
	fmt.Println()
}

func (cmd *ConfigShowCommand) Run(args []string) error {
	slog.Debug("ConfigShowCommand.Run()", "args", args)
	// Parse flags & get positional args
	err := cmd.HandleArgs(args)
	if err != nil { return err }
	// Execute command
	err = RunConfigShow()
	if err != nil { return err }

	return nil
}

func RunConfigShow() error {
	err := common.ShowConfig(os.Stdout)
	if err != nil { return err }
	
	return nil
}

// func CreateConfigCommand() *cobra.Command {
// 	command := cobra.Command{
// 		Use:   "config",
// 		Short: "config commands",
// 		Long:  "commands for getting, setting, etc config values",
// 	}
// 	command.AddCommand(CreateConfigGetCommand())
// 	command.AddCommand(CreateConfigSetCommand())
// 	command.AddCommand(CreateConfigShowCommand())
// 	return &command
// }

// func CreateConfigGetCommand() *cobra.Command {
// 	command := cobra.Command{
// 		Use:   "get field",
// 		Short: "Get config value",
// 		Long:  "Get individual config value from config settings",
// 		RunE:  runConfigGetCommand,
// 	}
// 	return &command
// }

// func runConfigGetCommand(cmd *cobra.Command, args []string) error {
// 	// Load config
// 	cfg, err := config.LoadConfig()
// 	if err != nil { return cfg, err }
// 	field := args[0]
// 	switch field {
// 	case "etcd-domain":
// 		domain, err := config.GetEtcDomain()
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Println(domain)
// 	default:
// 		errmsg := fmt.Sprintf("Unknown field: %s", field)
// 		return errors.New(errmsg)
// 	}
// 	return nil
// }

// func CreateConfigSetCommand() *cobra.Command {
// 	command := cobra.Command{
// 		Use:   "set key value",
// 		Short: "Set a config value",
// 		Long:  "Set a config value",
// 	}
// 	command.AddCommand(CreateConfigSetDomainCommand())
// 	return &command
// }

// func CreateConfigSetDomainCommand() *cobra.Command {
// 	command := cobra.Command{
// 		Use:   "etcd-domain $etcd-domain",
// 		Short: "Set the etcd domain",
// 		Long:  "Set the etcd domain",
// 		RunE:  runConfigSetDomainCommand,
// 		Args:  cobra.ExactArgs(1),
// 	}
// 	return &command
// }

// func runConfigSetDomainCommand(cmd *cobra.Command, args []string) error {
// 	etcdDomain := args[0]
// 	config, err := lib.LoadConfig()
// 	if err != nil {
// 		return err
// 	}
// 	config.EtcdDomain = etcdDomain
// 	err = lib.SaveConfig(config)
// 	return err
// }

// func CreateConfigShowCommand() *cobra.Command {
// 	command := cobra.Command{
// 		Use:   "show",
// 		Short: "show the current config",
// 		Long:  "show the current config",
// 		RunE:  runConfigShowCommand,
// 		Args:  cobra.ExactArgs(0),
// 	}
// 	return &command
// }

// func runConfigShowCommand(cmd *cobra.Command, args []string) error {
// 	err := lib.ShowConfig(os.Stdout)
// 	return err
// }
