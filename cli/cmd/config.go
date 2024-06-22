package cmd
import (
	
	"github.com/spf13/cobra"
	
	dylt "github.com/dylt-dev/dylt/lib"
)

func CreateConfigCommand () *cobra.Command {
	command := cobra.Command {
		Use: "config",
		Short: "config commands",
		Long: "commands for getting, setting, etc config values",
	}
	return &command
}


func CreateConfigGetCommand () *cobra.Command {
	command := cobra.Command {
		Use: "config get",
		Short: "Get config value",
		Long: "Get individual config value from config settings",
		RunE: runConfigGetCommand,
	}
	return &command
}

func runConfigGetCommand (cmd *cobra.Command, args []string) error {
	
	return nil
}