package cmd

import (
	"flag"
	"fmt"
	"log/slog"
)

type ConfigCommand struct {
	*flag.FlagSet
}

func NewConfigCommand () *ConfigCommand {
	flagSet := flag.NewFlagSet("config", flag.PanicOnError)
	return &ConfigCommand { FlagSet: flagSet }
}


func createConfigSubCommand (cmdName string) (Command, error) {
	switch cmdName {
	case "get": return NewConfigGetCommand(), nil
	case "set": return NewConfigSetCommand(), nil
	case "show": return NewConfigShowCommand(), nil
	default: return nil, fmt.Errorf("unrecognized command: %s", cmdName)
	}
}

func (cmd *ConfigCommand) Run (args []string) error {
	slog.Debug("ConfigCommand.Run()", "args", args)
	// Pargs flags & get positional params
	err := cmd.Parse(args)
	if err != nil { return err }
	var cmdline Cmdline = cmd.Args()
	// Get the subcommand
	if !cmdline.HasCommand() {
		return fmt.Errorf("config requires subcommand")
	}
	subCmdName := cmdline.Command()
	subCmd, err := createConfigSubCommand(subCmdName)
	if err != nil { return err }
	subArgs := cmdline.Args()
	slog.Debug("ConfigCommand.Run()", "subCmdName", subCmdName, "subArgs", subArgs)
	err = subCmd.Run(subArgs)
	if err != nil { return err }
	return nil
}

type ConfigGetCommand struct {
	*flag.FlagSet
}

func (cmd *ConfigGetCommand) Run (args []string) error {
	slog.Debug("ConfigGetCommand.Run()", "args", args)
	err := cmd.Parse(args)
	if err != nil { return err }
	return nil
}

func NewConfigGetCommand () *ConfigGetCommand {
	flagSet := flag.NewFlagSet("config.get", flag.PanicOnError)
	return &ConfigGetCommand { FlagSet: flagSet }
}

type ConfigSetCommand struct {
	*flag.FlagSet
}

func NewConfigSetCommand () *ConfigSetCommand {
	flagSet := flag.NewFlagSet("config.set", flag.PanicOnError)
	return &ConfigSetCommand { FlagSet: flagSet }
}

func (cmd *ConfigSetCommand) Run (args []string) error {
	slog.Debug("ConfigSetCommand.Run()", "args", args)
	err := cmd.Parse(args)
	if err != nil { return err }
	return nil
}

type ConfigShowCommand struct {
	*flag.FlagSet
}

func NewConfigShowCommand () *ConfigShowCommand {
	flagSet := flag.NewFlagSet("config.show", flag.PanicOnError)
	return &ConfigShowCommand { FlagSet: flagSet }
}

func (cmd *ConfigShowCommand) Run (args []string) error {
	slog.Debug("ConfigShowCommand.Run()", "args", args)
	err := cmd.Parse(args)
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
// 		Use:   "etcd_domain $etcd_domain",
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
