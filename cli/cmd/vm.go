package cmd

import (
	// "encoding/json"

	"encoding/json"
	"flag"
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/lib"
)

type VmCommand struct {
	*flag.FlagSet
}

type VmAddCommand struct {
	*flag.FlagSet
}

func NewVmAddCommand () *VmAddCommand {
	flagSet := flag.NewFlagSet("vm.add", flag.PanicOnError)
	return &VmAddCommand{FlagSet: flagSet}
}

func (cmd VmAddCommand) Run (args []string) error {
	// parse flags + setup positional args
	slog.Debug("VmAddCommand.Run()", "args", args)
	err := cmd.Parse(args)
	if err != nil { return err }
	cmdArgs := cmd.Args()
	if len(cmdArgs) != 2 { return fmt.Errorf("vm add requires 2 args; found %d", len(cmdArgs))}
	// get vm-specific etcd client
	cli, err := lib.CreateVmClientFromConfig()
	if err != nil { return err }
	// execute command
	name := args[0]
	address := args[1]
	vm, err := cli.Add(name, address)
	if err != nil { return err }
	fmt.Println(vm)
	return nil
}



type VmAllCommand struct {
	*flag.FlagSet
}

func NewVmAllCommand () *VmAllCommand {
	flagSet := flag.NewFlagSet("vm.all", flag.PanicOnError)
	return &VmAllCommand{FlagSet: flagSet}
}

func (cmd VmAllCommand) Run (args[] string) error {
	// parse flags + setup positional args
	slog.Debug("VmAllCommand.Run()", "args", args)
	err := cmd.Parse(args)
	if err != nil { return err }
	cmdArgs := cmd.Args()
	if len(cmdArgs) != 0 { return fmt.Errorf("vm all requires 0 arguments; %d found", len(cmdArgs))}
	// get vm-specific etcd client
	cli, err := lib.CreateVmClientFromConfig()
	if err != nil { return err }
	// execute command
 	vmInfoMap, err := cli.All()
	if err != nil { return err }
	jsonData, err := json.Marshal(vmInfoMap)
	if err != nil { return err }
	fmt.Println(string(jsonData))
	return nil
}

type VmDelCommand struct {
	*flag.FlagSet
}

func (cmd *VmDelCommand) ParseAndValidate (args []string, nArgs int) ([]string, error) {
	err := cmd.Parse(args)
	if err != nil { return nil, err }
	cmdArgs := cmd.Args()
	if len(cmdArgs) != nArgs { return nil, fmt.Errorf("vm all requires %d argument(s); %d found", nArgs, len(cmdArgs))}
	return cmdArgs, nil
}

func NewVmDelCommand () *VmDelCommand {
	flagSet := flag.NewFlagSet("vm.del", flag.PanicOnError)
	return &VmDelCommand{FlagSet: flagSet}
}

func (cmd VmDelCommand) Run (args[] string) error {
	// parse flags and setup positional args
	slog.Debug("VmDelCommand.Run()", "args", args)
	cmdArgs, err := cmd.ParseAndValidate(args, 1)
	if err != nil { return err }
	name := cmdArgs[0]
	// get vm-specific etcd client
	cli, err := lib.CreateVmClientFromConfig()
	if err != nil { return err }
	// execute command
	prevVal, err := cli.Del(name)
	if err != nil { return err }
	fmt.Printf("%s\n", string(prevVal))
	return nil
}

type VmGetCommand struct {
	*flag.FlagSet
}

func NewVmGetCommand () *VmGetCommand {
	flagSet := flag.NewFlagSet("vm.get", flag.PanicOnError)
	return &VmGetCommand{FlagSet: flagSet}
}

func ParseAndValidate (flagSet *flag.FlagSet, args []string, nArgs int) ([]string, error) {
	err := flagSet.Parse(args)
	if err != nil { return nil, err }
	cmdArgs := flagSet.Args()
	if len(cmdArgs) != nArgs { return nil, fmt.Errorf("vm all requires %d argument(s); %d found", nArgs, len(cmdArgs))}
	return cmdArgs, nil
}

func (cmd *VmGetCommand) Run (args[] string) error {
	// parse flags and setup positional args
	slog.Debug("VmGetCommand.Run()", "args", args)
	cmdArgs, err := ParseAndValidate(cmd.FlagSet, args, 1)
	if err != nil { return err }
	name := cmdArgs[0]
	// get vm-specific etcd client
	cli, err := lib.CreateVmClientFromConfig()
	if err != nil { return err }
	// execute command
	vm, err := cli.Get(name)
	if err != nil { return err }
	fmt.Println(vm)
	return nil

}

type VmListCommand struct {
	*flag.FlagSet
}

func NewVmListCommand () *VmListCommand {
	flagSet := flag.NewFlagSet("vm.list", flag.PanicOnError)
	return &VmListCommand{FlagSet: flagSet}
}

func (cmd VmListCommand) Run (args[] string) error {
	// parse flags and setup positional args
	slog.Debug("VmListCommand.Run()", "args", args)
	_, err := ParseAndValidate(cmd.FlagSet, args, 0)
	if err != nil { return err }
	// get vm-specific etcd client
	cli, err := lib.CreateVmClientFromConfig()
	if err != nil { return err }
	// execute command
	names, err := cli.Names()
	if err != nil { return err }
	jsonData, err := json.Marshal(names)
	if err != nil { return err }
	fmt.Println(string(jsonData))
	return nil
}

type VmSetCommand struct {
	*flag.FlagSet
}

func NewVmSetCommand () *VmSetCommand {
	flagSet := flag.NewFlagSet("vm.set", flag.PanicOnError)
	return &VmSetCommand{FlagSet: flagSet}
}

func (cmd VmSetCommand) Run (args[] string) error {
	// parse flags and setup positional args
	slog.Debug("VmSetCommand.Run()", "args", args)
	cmdArgs, err := ParseAndValidate(cmd.FlagSet, args, 3)
	if err != nil { return err }
	name := cmdArgs[0]
	key := cmdArgs[1]
	val := cmdArgs[2]
	// get vm-specific etcd client
	cli, err := lib.CreateVmClientFromConfig()
	if err != nil { return err }
	// execute command
	vm, err := cli.Get(name)
	if err != nil { return err }
	err = vm.Set(key, val)
	if err != nil { return err }
	vm, err = cli.Put(name, vm)
	if err != nil { return err }
	fmt.Println(vm)
	return nil
}

type VmShowCommand struct {
	*flag.FlagSet
}

func NewVmShowCommand () *VmShowCommand {
	flagSet := flag.NewFlagSet("vm.show", flag.PanicOnError)
	return &VmShowCommand{FlagSet: flagSet}
}

func (cmd VmShowCommand) Run (args[] string) error {
	slog.Debug("VmShowCommand.Run()", "args", args)
	err := cmd.Parse(args)
	if err != nil { return err }
	return nil
}


func createVmSubCommand (cmd string) (Command, error) {
	switch cmd {
	case "add": return NewVmAddCommand(), nil
	case "all": return NewVmAllCommand(), nil
	case "del": return NewVmDelCommand(), nil
	case "get": return NewVmGetCommand(), nil
	case "list": return NewVmListCommand(), nil
	case "set": return NewVmSetCommand(), nil
	// case "show": return NewVmShowCommand(), nil
	default: return nil, fmt.Errorf("unrecognized subcommand: %s", cmd)
	}
}

func NewVmCommand () *VmCommand {
	flagSet := flag.NewFlagSet("vm", flag.PanicOnError)
	return &VmCommand{FlagSet: flagSet}
}


func (cmd *VmCommand) Run (args []string) error {
	slog.Debug("VmCommand.Run()", "args", args)
	// Parse flags & get positional params
	err := cmd.Parse(args)
	if err != nil { return err }
	var cmdline Cmdline = cmd.Args()
	// Get the subcommand
	if !cmdline.HasCommand() {
		return fmt.Errorf("vm requires subcommand")
	}
	subCmdName := cmdline.Command()
	subCmd, err := createVmSubCommand(subCmdName)
	if err != nil { return err }
	subArgs := cmdline.Args()
	err = subCmd.Run(subArgs)
	if err != nil { return err }
	return nil
}

// func CreateVmCommand () *cobra.Command {
// 	command := cobra.Command {
// 		Use: "vm",
// 		Short: "VM commands",
// 		Long: "Operations on VM objects in storage",
// 	}
// 	command.AddCommand(CreateVmAddCommand())
// 	command.AddCommand(CreateVmAllCommand())
// 	command.AddCommand(CreateVmDelCommand())
// 	command.AddCommand(CreateVmGetCommand())
// 	command.AddCommand(CreateVmListCommand())
// 	command.AddCommand(CreateVmSetCommand())
// 	command.AddCommand(CreateVmShowCommand())
// 	return &command
// }

// func CreateVmAddCommand () *cobra.Command {
// 	command := cobra.Command {
// 		Use: "add name address",
// 		Short: "Add new VM",
// 		Long: "Add a new VM to the collection",
// 		RunE: runVmAddCommand,
// 		Args: cobra.ExactArgs(2),
// 	}
// 	return &command
// }

// func runVmAddCommand (cmd *cobra.Command, args []string) error {
// 	cli, err := dylt.CreateVmClientFromConfig()
// 	if err != nil { return err }
// 	name := args[0]
// 	address := args[1]
// 	vm, err := cli.Add(name, address)
// 	if err != nil { return err }
// 	fmt.Println(vm)

// 	return nil
// }

// func CreateVmAllCommand () *cobra.Command {
// 	command := cobra.Command {
// 		Use: "all",
// 		Short: "All VM info",
// 		Long: "Return data for all VMs in the system",
// 		RunE: runVmAllCommand,
// 	}
// 	command.Flags().BoolP("just-names", "n", false, "return just the host names")
// 	return &command
// }

// func runVmAllCommand (cmd *cobra.Command, args []string) error {
// 	cli, err := dylt.CreateVmClientFromConfig()
// 	if err != nil { return err }
// 	vmInfoMap, err := cli.All()
// 	if err != nil { return err }
// 	// hasShr := cmd.Flags().Changed("shr")
// 	// if hasShr {
// 	// 	shr, err := cmd.Flags().GetBool("shr")
// 	// 	if err != nil { return err }
// 	// 	vmInfoMap = dylt.FilterOnShr(vmInfoMap, shr)
// 	// }
// 	isJustNames, err := cmd.Flags().GetBool("just-names")
// 	if err != nil { return err }
// 	if isJustNames {
// 		names := []string{}
// 		for name, _ := range(vmInfoMap) {
// 			names = append(names, name)
// 		}
// 		buf, err := json.Marshal(names)
// 		if err != nil { return err }
// 		fmt.Println(string(buf))
// 	} else {
// 		jsonData, err := json.Marshal(vmInfoMap)
// 		if err != nil { return err }
// 		fmt.Println(string(jsonData))
// 	}
// 	return nil
// }

// func CreateVmDelCommand () *cobra.Command {
// 	command := cobra.Command {
// 		Use: "del",
// 		Short: "Delete a VM entry",
// 		Long: "Delete a VM entry from etcd",
// 		RunE: runVmDelCommand,
// 	}
// 	return &command
// }

// func runVmDelCommand (cmd *cobra.Command, args []string) error {
// 	cli, err := dylt.CreateVmClientFromConfig()
// 	if err != nil { return err }
// 	arg := args[0]
// 	key := fmt.Sprintf("/vm/%s", arg)
// 	prevVal, err := cli.Delete(key)
// 	if err != nil { return err }
// 	if prevVal == nil { return nil }
// 	fmt.Println(string(prevVal))
// 	return nil
// }

// func CreateVmGetCommand () *cobra.Command {
// 	command := cobra.Command {
// 		Use: "get $vm",
// 		Short: "Get a VM",
// 		Long: "Get information on a VM",
// 		RunE: runVmGetCommand,
// 		Args: cobra.RangeArgs(1, 2),
// 	}
// 	return &command
// }

// func runVmGetCommand (cmd *cobra.Command, args []string) error {
// 	key := args[0]
// 	attr := ""
// 	if len(args) >= 2 {
// 		attr = args[1]
// 	}
// 	hasAttr := false
// 	if attr != "" {
// 		hasAttr = true
// 	} 
// 	cli, err := dylt.CreateVmClientFromConfig()
// 	if err != nil { return err }
// 	vm, err := cli.Get(key)
// 	if err != nil { return err }
// 	if vm == nil { return nil }
// 	if hasAttr {
// 		attrValue, err := vm.Get(attr)
// 		if err != nil { return err }
// 		fmt.Printf("%s\n", attrValue)
// 	} else {
// 		fmt.Println(vm)
// 	}
// 	return nil
// }

// func CreateVmListCommand () *cobra.Command {
// 	command := cobra.Command {
// 		Use: "list",
// 		Short: "List VMs",
// 		Long: "List all VMs in the system",
// 		RunE: runVmListCommand,
// 	}
// 	return &command
// }

// func runVmListCommand (cmd *cobra.Command, args []string) error {
// 	cli, err := dylt.CreateVmClientFromConfig()
// 	if err != nil { return err }
// 	names, err := cli.Names()
// 	if err != nil { return err }
// 	jsonData, err := json.Marshal(names)
// 	if err != nil { return err }
// 	fmt.Println(string(jsonData))
// 	return nil
// }

// func CreateVmSetCommand () *cobra.Command {
// 	command := cobra.Command {
//  		Use: "set vm key value",
// 		Short: "Set a VM attribute",
// 		Long: "Set an attribute on a VM. Create the attribute if it doesn't exist.",
// 		RunE: runVmSetCommand,
// 		Args: cobra.ExactArgs(3),
// 	}
// 	return &command
// }

// func runVmSetCommand (cmd *cobra.Command, args []string) error {
// 	cli, err := dylt.CreateVmClientFromConfig()
// 	if err != nil { return err }
// 	defer cli.Close()
// 	name := args[0]
// 	field := args[1]
// 	value := args[2]
// 	key := fmt.Sprintf("/vm/%s", name)
// 	vm, err := cli.Get(name)
// 	if err != nil { return err }
// 	fmt.Printf("vm=%s\n", vm)
// 	switch field {
// 	case "Address":
// 		vm.Address = value
// 	default:
// 		errmsg := fmt.Sprintf("Unknown field: %s", field)
// 		err := errors.New(errmsg)
// 		return err
// 	}
// 	s, err := json.Marshal(vm)
// 	fmt.Printf("string(s)=%s\n", string(s))
// 	if err != nil { return err }
// 	ctx := context.Background()
// 	ctx, cancel := context.WithTimeout(ctx, 10 * time.Second)
// 	resp, err := cli.Client.Put(ctx, key, string(s))
// 	fmt.Printf("err=%s\n", err)
// 	fmt.Printf("resp=%v\n", resp)
// 	cancel()
// 	vmNew, err := cli.Get(name)
// 	if err != nil { return err }
// 	bufNew, err := json.Marshal(vmNew)
// 	if err != nil { return err }
// 	fmt.Println(string(bufNew))
// 	return nil
// }

// func CreateVmShowCommand () *cobra.Command {
// 	command := cobra.Command {
// 		Use: "show $vm",
// 		Short: "show a VM",
// 		Long: "show a VM",
// 		RunE: runVmShowCommand,
// 		Args: cobra.ExactArgs(1),
// 	}
// 	return &command
// }

// func runVmShowCommand (cmd *cobra.Command, args []string) error {
// 	vmName := args[0]
// 	vmClient, err := dylt.CreateVmClientFromConfig()
// 	if err != nil { return err }
// 	vm, err := vmClient.Get(vmName)
// 	if err != nil { return err }
// 	fmt.Println(vm)
// 	return nil
// }
