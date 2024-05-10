package cmd

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"log/slog"

	// "os"

	"github.com/spf13/cobra"

	dylt "github.com/dylt-dev/dylt/lib"
)

func CreateVmCommand () *cobra.Command {
	command := cobra.Command {
		Use: "vm",
		Short: "VM commands",
		Long: "Operations on VM objects in storage",
	}
	command.AddCommand(CreateVmAllCommand())
	command.AddCommand(CreateVmGetCommand())
	command.AddCommand(CreateVmListCommand())
	return &command
}


func CreateVmAllCommand () *cobra.Command {
	command := cobra.Command {
		Use: "all",
		Short: "All VM info",
		Long: "Return data for all VMs in the system",
		RunE: runVmAllCommand,
	}
	return &command
}

func runVmAllCommand (cmd *cobra.Command, args []string) error {
	cli, err := dylt.NewVmClient("hello.dylt.dev")
	if err != nil {
		slog.Error("Error creating new vm client")
		return err
	}
	names, err := cli.Names()
	if err != nil { return err }
	vms := []*dylt.VmInfo{}
	for _, name := range(names) {
		vmInfo, err := cli.Get(name)
		if err != nil { return err }
		vms = append(vms, vmInfo)
	}
	jsonData, err := json.Marshal(vms)
	if err != nil { return err }
	fmt.Println(string(jsonData))
	return nil
}


func CreateVmGetCommand () *cobra.Command {
	command := cobra.Command {
		Use: "get",
		Short: "Get a VM",
		Long: "Get information on a VM",
		RunE: runVmGetCommand,
	}
	return &command
}

func runVmGetCommand (cmd *cobra.Command, args []string) error {
	cli, err := dylt.NewVmClient("hello.dylt.dev")
	if err != nil {
		slog.Error("Error creating new vm client")
		return err
	}
	key := args[0]
	vm, err := cli.Get(key)
	if err != nil { return err }
	jsonData, err := json.Marshal(vm)
	if err != nil { return err }
	fmt.Println(string(jsonData))
	return nil
}

func CreateVmListCommand () *cobra.Command {
	command := cobra.Command {
		Use: "list",
		Short: "List VMs",
		Long: "List all VMs in the system",
		RunE: runVmListCommand,
	}
	return &command
}

func runVmListCommand (cmd *cobra.Command, args []string) error {
	cli, err := dylt.NewVmClient("hello.dylt.dev")
	if err != nil {
		slog.Error("Error creating new vm client")
		return err
	}
	names, err := cli.Names()
	if err != nil { return err }
	jsonData, err := json.Marshal(names)
	if err != nil { return err }
	fmt.Println(string(jsonData))
	return nil
}