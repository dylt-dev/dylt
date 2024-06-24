package cmd

import (
	"errors"
	"fmt"

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