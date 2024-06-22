package cmd

import (
	// "encoding/json"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	dylt "github.com/dylt-dev/dylt/lib"
)

func CreateVmCommand () *cobra.Command {
	command := cobra.Command {
		Use: "vm",
		Short: "VM commands",
		Long: "Operations on VM objects in storage",
	}
	command.AddCommand(CreateVmAddCommand())
	command.AddCommand(CreateVmAllCommand())
	command.AddCommand(CreateVmDelCommand())
	command.AddCommand(CreateVmGetCommand())
	command.AddCommand(CreateVmListCommand())
	command.AddCommand(CreateVmSetCommand())
	return &command
}


func CreateVmAddCommand () *cobra.Command {
	command := cobra.Command {
		Use: "add name address",
		Short: "Add new VM",
		Long: "Add a new VM to the collection",
		RunE: runVmAddCommand,
		Args: cobra.ExactArgs(2),
	}
	return &command
}

func runVmAddCommand (cmd *cobra.Command, args []string) error {
	cli, err := dylt.NewVmClient(viper.GetString("etcd_domain"))
	if err != nil { return err }
	name := args[0]
	address := args[1]
	vm, err := cli.Add(name, address)
	if err != nil { return err }
	fmt.Println(vm)

	return nil
}


func CreateVmAllCommand () *cobra.Command {
	command := cobra.Command {
		Use: "all",
		Short: "All VM info",
		Long: "Return data for all VMs in the system",
		RunE: runVmAllCommand,
	}
	command.Flags().Bool("shr", true, "return hosts with (or without) SHR")
	return &command
}


func runVmAllCommand (cmd *cobra.Command, args []string) error {
	cli, err := dylt.NewVmClient(viper.GetString("etcd_domain"))
	if err != nil { return err }
	vmInfoMap, err := cli.All()
	if err != nil { return err }
	// hasShr := cmd.Flags().Changed("shr")
	// if hasShr {
	// 	shr, err := cmd.Flags().GetBool("shr")
	// 	if err != nil { return err }
	// 	vmInfoMap = dylt.FilterOnShr(vmInfoMap, shr)
	// }
	jsonData, err := json.Marshal(vmInfoMap)
	if err != nil { return err }
	fmt.Println(string(jsonData))
	return nil
}


func CreateVmDelCommand () *cobra.Command {
	command := cobra.Command {
		Use: "del",
		Short: "Delete a VM entry",
		Long: "Delete a VM entry from etcd",
		RunE: runVmDelCommand,
	}
	return &command
}

func runVmDelCommand (cmd *cobra.Command, args []string) error {
	cli, err := dylt.NewVmClient(viper.GetString("etcd_domain"))
	if err != nil { return err }
	arg := args[0]
	key := fmt.Sprintf("/vm/%s", arg)
	prevVal, err := cli.Delete(key)
	if err != nil { return err }
	if prevVal == nil { return nil }
	fmt.Println(string(prevVal))
	return nil
}


func CreateVmGetCommand () *cobra.Command {
	command := cobra.Command {
		Use: "get $vm",
		Short: "Get a VM",
		Long: "Get information on a VM",
		RunE: runVmGetCommand,
		Args: cobra.RangeArgs(1, 2),
	}
	return &command
}

func runVmGetCommand (cmd *cobra.Command, args []string) error {
	key := args[0]
	attr := ""
	if len(args) >= 2 {
		attr = args[1]
	}
	hasAttr := false
	if attr != "" {
		hasAttr = true
	} 
	cli, err := dylt.NewVmClient(viper.GetString("etcd_domain"))
	if err != nil { return err }
	vm, err := cli.Get(key)
	if err != nil { return err }
	if vm == nil { return nil }
	if hasAttr {
		attrValue, err := vm.Get(attr)
		if err != nil { return err }
		fmt.Printf("%s\n", attrValue)
	} else {
		fmt.Println(vm)
	}
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
	cli, err := dylt.NewVmClient(viper.GetString("etcd_domain"))
	if err != nil { return err }
	names, err := cli.Names()
	if err != nil { return err }
	jsonData, err := json.Marshal(names)
	if err != nil { return err }
	fmt.Println(string(jsonData))
	return nil
}

func CreateVmSetCommand () *cobra.Command {
	command := cobra.Command {
 		Use: "set vm key value",
		Short: "Set a VM attribute",
		Long: "Set an attribute on a VM. Create the attribute if it doesn't exist.",
		RunE: runVmSetCommand,
		Args: cobra.ExactArgs(3),
	}
	return &command
}

func runVmSetCommand (cmd *cobra.Command, args []string) error {
	cli, err := dylt.NewVmClient(viper.GetString("etcd_domain"))
	if err != nil { return err }
	defer cli.Close()
	name := args[0]
	field := args[1]
	value := args[2]
	key := fmt.Sprintf("/vm/%s", name)
	vm, err := cli.Get(name)
	if err != nil { return err }
	fmt.Printf("vm=%s\n", vm)
	switch field {
	case "Address":
		vm.Address = value
	default:
		errmsg := fmt.Sprintf("Unknown field: %s", field)
		err := errors.New(errmsg)
		return err
	}
	s, err := json.Marshal(vm)
	fmt.Printf("string(s)=%s\n", string(s))
	if err != nil { return err }
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10 * time.Second)
	resp, err := cli.Client.Put(ctx, key, string(s))
	fmt.Printf("err=%s\n", err)
	fmt.Printf("resp=%v\n", resp)
	cancel()
	vmNew, err := cli.Get(name)
	if err != nil { return err }
	bufNew, err := json.Marshal(vmNew)
	if err != nil { return err }
	fmt.Println(string(bufNew))
	return nil
}
