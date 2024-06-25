package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	dylt "github.com/dylt-dev/dylt/lib"
)

func CreateConfigCommand () *cobra.Command {
	command := cobra.Command {
		Use: "config",
		Short: "config commands",
		Long: "commands for getting, setting, etc config values",
	}
	command.AddCommand(CreateConfigGetCommand())
	command.AddCommand(CreateConfigSetCommand())
	command.AddCommand(CreateConfigShowCommand())
	return &command
}


func CreateConfigGetCommand () *cobra.Command {
	command := cobra.Command {
		Use: "get field",
		Short: "Get config value",
		Long: "Get individual config value from config settings",
		RunE: runConfigGetCommand,
	}
	return &command
}

func runConfigGetCommand (cmd *cobra.Command, args []string) error {
	config := dylt.Config{}
	err := config.Load()
	if err != nil { return err }
	field := args[0]
	switch field {
	case "etcd_domain":
		domain, err := config.GetEtcDomain()
		if err != nil { return err }
		fmt.Println(domain)
	default:
		errmsg := fmt.Sprintf("Unknown field: %s", field)
		return errors.New(errmsg)
	}
	return nil	
}

func CreateConfigSetCommand () *cobra.Command {
	command := cobra.Command {
		Use: "set key value",
		Short: "Set a config value",
		Long: "Set a config value",
	}
	command.AddCommand(CreateConfigSetDomainCommand())
	return &command
}

func CreateConfigSetDomainCommand () *cobra.Command {
	command := cobra.Command {
		Use: "etcd_domain $etcd_domain",
		Short: "Set the etcd domain",
		Long: "Set the etcd domain",
		RunE: runConfigSetDomainCommand,
		Args: cobra.ExactArgs(1),
	}
	return &command
}

func runConfigSetDomainCommand (cmd *cobra.Command, args []string) error {
	etcdDomain := args[0]
	config, err := dylt.LoadConfig()
	if err != nil { return err }
	err = config.SetEtcDomainAndSave(etcdDomain)
	return err
}

func CreateConfigShowCommand () *cobra.Command {
	command := cobra.Command {
		Use: "show",
		Short: "show the current config",
		Long: "show the current config",
		RunE: runConfigShowCommand,
		Args: cobra.ExactArgs(0),
	}
	return &command
}

func runConfigShowCommand (cmd *cobra.Command, args []string) error {
	err := dylt.ShowConfig(os.Stdout)
	return err
}
