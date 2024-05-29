package cmd
import (
	"fmt"
	
	"github.com/spf13/cobra"
	
	dylt "github.com/dylt-dev/dylt/lib"
)

func CreateCallCommand () *cobra.Command {
	command := cobra.Command {
		Use: "call [args]",
		Short: "call",
		Long: "call",
		RunE: runCallCommand,
	}
	return &command
}

func runCallCommand (cmd *cobra.Command, args []string) error {
	_, stdout, err := dylt.CallDaylightScriptO(args)
	if err != nil { return err }
	fmt.Printf("%s\n", stdout)
	return nil
}