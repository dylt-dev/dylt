package cmd
import (
	
	"github.com/spf13/cobra"
	
	dylt "github.com/dylt-dev/dylt/lib"
)

func CreateInitCommand () *cobra.Command {
	command := cobra.Command {
		Use: "init",
		Short: "Initialize dylt",
		Long: "Initialize dylt",
		RunE: runInitCommand,
	}
	command.Flags().String("etcd-domain", "", "etcd cluster to activate")
	command.MarkFlagRequired("etcd-domain")
	return &command
}

func runInitCommand (cmd *cobra.Command, args []string) error {
	etcdDomain, err := cmd.Flags().GetString("etcd-domain")
	if err != nil { return err }
	initInfo := dylt.InitInfo{
		EtcdDomain: etcdDomain,
	}
	err = dylt.Init(&initInfo)
	if err != nil { return err }
	return nil
}