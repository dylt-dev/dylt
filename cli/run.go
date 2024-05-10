package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/dylt-dev/dylt/cli/cmd"
)

func Run () int {
	rootCmd := createRootCommand()
	rootCmd.AddCommand(cmd.CreateGetCommand())
	rootCmd.AddCommand(cmd.CreateVmCommand())
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func createRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "dylt",
		Short: "dylt core functions",
		Long: "CLI for using core daylight (dylt) features",
	}
	return cmd
}