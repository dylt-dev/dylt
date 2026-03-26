package cmd

import (
	// "encoding/json"

	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/common"
	"github.com/dylt-dev/dylt/eco"
)

type VmCommand struct {
	*BaseCommand
}

func NewVmCommand(cmdline Cmdline, parent Command) *VmCommand {
	// vm command
	name := "vm"
	cmd := &VmCommand{BaseCommand: NewBaseCommand(name, cmdline, parent)}
	
	//init flags (if any)
	
	return cmd
}

func (cmd *VmCommand) CreateSubCommand() (Command, error) {
	args, flag := cmd.Args()
	if !flag {
		return nil, nil
	}
	return createVmSubCommand(args, cmd)
}

func (cmd *VmCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}

	// if Help flag is set, no further processing is necessary
	if cmd.Help {
		return nil
	}

	// validate arg count
	cmdArgs, _ := cmd.Args()
	cmdName := "vm"
	nExpected := 0
	if len(cmdArgs) < nExpected {
		return fmt.Errorf("`%s` expects >=%d argument(s); received %d",
			cmdName,
			nExpected,
			len(cmdArgs))
	}

	return nil
}

func (cmd *VmCommand) PrintUsage() {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Vm_Add_Short)
	fmt.Printf("\t%s\n", USG_Vm_All_Short)
	fmt.Printf("\t%s\n", USG_Vm_Del_Short)
	fmt.Printf("\t%s\n", USG_Vm_Get_Short)
	fmt.Printf("\t%s\n", USG_Vm_List_Short)
	fmt.Printf("\t%s\n", USG_Vm_Set_Short)
	fmt.Println()
}

func (cmd *VmCommand) Run() error {
	slog.Debug("VmCommand.Run()", "args", cmd.Cmdline)

	// Check for 0 args; if so print usage & return
	if len(cmd.Cmdline) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}

	// If help flag set, print usage + return
	if cmd.Help {
		cmd.PrintUsage()
		return nil
	}

	// If no args, print usage
	args, _ := cmd.Args()
	if len(args) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// execute command
	err = RunVm(args, cmd)
	return err
}

func RunVm(cmdline Cmdline, parent *VmCommand) error {
	slog.Debug("RunVm()", "cmdline", cmdline, "parent", parent)
	// create the subcommand and run it
	subCmd, err := createVmSubCommand(cmdline, parent)
	if err != nil {
		return err
	}
	err = subCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func createVmSubCommand(cmdline Cmdline, parent Command) (Command, error) {
	subCmd := cmdline.Command()
	switch subCmd {
	case "add":
		return NewVmAddCommand(cmdline, parent), nil
	case "all":
		return NewVmAllCommand(cmdline, parent), nil
	case "del":
		return NewVmDelCommand(cmdline, parent), nil
	case "get":
		return NewVmGetCommand(cmdline, parent), nil
	case "list":
		return NewVmListCommand(cmdline, parent), nil
	case "set":
		return NewVmSetCommand(cmdline, parent), nil
	default:
		return nil, fmt.Errorf("unrecognized subcommand: %s", subCmd)
	}
}

type VmAddCommand struct {
	*BaseCommand
	Name string // arg 0
	Fqdn string // arg 1
}

func NewVmAddCommand(cmdline Cmdline, parent Command) *VmAddCommand {
	// vm add command
	name := "vm.add"
	cmd := &VmAddCommand{BaseCommand: NewBaseCommand(name, cmdline, parent)}
	
	//init flags (if any)
	
	return cmd

}

func (cmd *VmAddCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}

	// if Help flag is set, no further processing is necessary
	if cmd.Help {
		return nil
	}

	// validate arg count
	cmdArgs, _ := cmd.Args()
	nExpected := 2
	if len(cmdArgs) != nExpected {
		cmd.PrintUsage()
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
			cmdString,
			nExpected,
			len(cmdArgs))
	}
	// init positional params
	cmd.Name = cmdArgs[0]
	cmd.Fqdn = cmdArgs[1]

	return nil
}

func (cmd *VmAddCommand) PrintUsage() {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Vm_Add)
	fmt.Println()
}

func (cmd VmAddCommand) Run() error {
	slog.Debug("VmAddCommand.Run()", "args", cmd.Cmdline)

	// parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}

	// If help flag set, print usage + return
	if cmd.Help {
		cmd.PrintUsage()
		return nil
	}

	// execute command
	err = RunVmAdd(cmd.Name, cmd.Fqdn)
	return err
}

func RunVmAdd(name string, fqdn string) error {
	// get vm-specific etcd client
	cli, err := eco.CreateVmClientFromConfig()
	if err != nil {
		return err
	}
	// execute command
	vm, err := cli.Add(name, fqdn)
	if err != nil {
		return err
	}
	fmt.Println(vm)

	return nil
}

type VmAllCommand struct {
	*BaseCommand
}

func NewVmAllCommand(cmdline Cmdline, parent Command) *VmAllCommand {
	// vm all command
	name := "vm.all"
	cmd := &VmAllCommand{BaseCommand: NewBaseCommand(name, cmdline, parent)}
	
	//init flags (if any)
	
	return cmd
}

func (cmd *VmAllCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}

	// if Help flag is set, no further processing is necessary
	if cmd.Help {
		return nil
	}

	// validate arg count
	cmdArgs, _ := cmd.Args()
	nExpected := 0
	if len(cmdArgs) != nExpected {
		cmd.PrintUsage()
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
			cmdString,
			nExpected,
			len(cmdArgs))
	}

	// init positional params (nop - no params)

	return nil
}

func (cmd *VmAllCommand) PrintUsage() {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Vm_All)
	fmt.Println()
}

func (cmd VmAllCommand) Run() error {
	slog.Debug("VmAllCommand.Run()", "args", cmd.Cmdline)

	// Parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}

	// If help flag set, print usage + return
	if cmd.Help {
		cmd.PrintUsage()
		return nil
	}

	// Execute command
	err = RunVmAll()
	return err
}

func RunVmAll() error {
	// get vm-specific etcd client, get all vm data, + show it
	cli, err := eco.CreateVmClientFromConfig()
	if err != nil {
		return err
	}
	vmInfoMap, err := cli.All()
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(vmInfoMap)
	if err != nil {
		return err
	}
	fmt.Println(string(jsonData))

	return nil
}

// Usage
//
//	vm del vmName
type VmDelCommand struct {
	*BaseCommand
	Name string // arg 0
}

func NewVmDelCommand(cmdline Cmdline, parent Command) *VmDelCommand {
	// vm del command
	name := "vm.del"
	cmd := &VmDelCommand{BaseCommand: NewBaseCommand(name, cmdline, parent)}

	return cmd
}

func (cmd *VmDelCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}

	// if Help flag is set, no further processing is necessary
	if cmd.Help {
		return nil
	}

	// validate arg count
	cmdArgs, _ := cmd.Args()
	nExpected := 1
	if len(cmdArgs) != nExpected {
		cmd.PrintUsage()
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
			cmdString,
			nExpected,
			len(cmdArgs))
	}

	// init positional params (nop - no params)
	cmd.Name = cmdArgs[0]

	return nil
}

func (cmd *VmDelCommand) PrintUsage() {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Vm_Del)
	fmt.Println()
}

func (cmd *VmDelCommand) Run() error {
	slog.Debug("VmDelCommand.Run()", "args", cmd.Cmdline)

	// parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}

	// If help flag set, print usage + return
	if cmd.Help {
		cmd.PrintUsage()
		return nil
	}

	// execute command
	err = RunVmDel(cmd.Name)
	return err
}

func RunVmDel(name string) error {
	slog.Debug("RunVmDel()", "name", name)
	// get vm-specific etcd client
	cli, err := eco.CreateVmClientFromConfig()
	if err != nil {
		return err
	}
	// delete vm data from cluster
	prevVal, err := cli.Del(name)
	if err != nil {
		return err
	}
	// log deleted vm, if it existede
	if prevVal == nil {
		fmt.Printf("vm '%s' not found\n", name)
	} else {
		fmt.Printf("%s\n", string(prevVal))
	}

	return nil
}

// Usage
//
//	vm get vmName
type VmGetCommand struct {
	*BaseCommand
	Name string // arg 0
}

func NewVmGetCommand(cmdline Cmdline, parent Command) *VmGetCommand {
	// vm get command
	name := "vm.get"
	cmd := &VmGetCommand{BaseCommand: NewBaseCommand(name, cmdline, parent)}
	
	//init flags (if any)
	
	return cmd
}

func (cmd *VmGetCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}

	// if Help flag is set, no further processing is necessary
	if cmd.Help {
		return nil
	}

	// validate arg count
	cmdArgs, _ := cmd.Args()
	nExpected := 1
	if len(cmdArgs) != nExpected {
		cmd.PrintUsage()
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
			cmdString,
			nExpected,
			len(cmdArgs))
	}

	// init positional params
	cmd.Name = cmdArgs[0]

	return nil
}

func (cmd *VmGetCommand) PrintUsage() {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Vm_Get)
	fmt.Println()
}

func (cmd *VmGetCommand) Run() error {
	slog.Debug("VmGetCommand.Run()", "args", cmd.Cmdline)

	// parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}

	// If help flag set, print usage + return
	if cmd.Help {
		cmd.PrintUsage()
		return nil
	}

	// execute command
	err = RunVmGet(cmd.Name)
	return err
}

func RunVmGet(name string) error {
	slog.Debug("RunVmGet()", "name", name)
	// get vm-specific etcd client
	cli, err := eco.CreateVmClientFromConfig()
	if err != nil {
		return err
	}
	// get vm data from cluster
	vm, err := cli.Get(name)
	if err != nil {
		return err
	}
	// pritn vm data if vm was found
	if vm == nil {
		fmt.Printf("\nvm '%s' not found.\n\n", name)
	} else {
		fmt.Println(vm)
	}

	return nil
}

// Usage
//
//	vm list
type VmListCommand struct {
	*BaseCommand
}

func NewVmListCommand(cmdline Cmdline, parent Command) *VmListCommand {
	// vm list command
	name := "vm.list"
	cmd := &VmListCommand{BaseCommand: NewBaseCommand(name, cmdline, parent)}
	
	//init flags (if any)
	
	return cmd
}

func (cmd *VmListCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}

	// if Help flag is set, no further processing is necessary
	if cmd.Help {
		return nil
	}

	// validate arg count
	cmdArgs, _ := cmd.Args()
	nExpected := 0
	if len(cmdArgs) != nExpected {
		cmd.PrintUsage()
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
			cmdString,
			nExpected,
			len(cmdArgs))
	}

	// init positional params (nop - no params)

	return nil
}

func (cmd *VmListCommand) PrintUsage() {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Vm_List)
	fmt.Println()
}

func (cmd VmListCommand) Run() error {
	slog.Debug("VmListCommand.Run()", "args", cmd.Cmdline)

	// parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}

	// If help flag set, print usage + return
	if cmd.Help {
		cmd.PrintUsage()
		return nil
	}

	// execute command
	err = RunVmList()
	return err
}

func RunVmList() error {
	slog.Debug("RunVmList()")
	// get vm-specific etcd client
	cli, err := eco.CreateVmClientFromConfig()
	if err != nil {
		return err
	}
	// List all vm names, one per line
	names, err := cli.Names()
	if err != nil {
		return err
	}
	fmt.Println()
	for _, name := range names {
		fmt.Println(name)
	}
	fmt.Println()

	return nil
}

// Usage
//
//	vm set vmName key val
type VmSetCommand struct {
	*BaseCommand
	Name  string // arg 0
	Key   string // arg 1
	Value string // arg 2
}

func NewVmSetCommand(cmdline Cmdline, parent Command) *VmSetCommand {
	// vm set command
	name := "vm.set"
	cmd := &VmSetCommand{BaseCommand: NewBaseCommand(name, cmdline, parent)}
	
	//init flags (if any)
	
	return cmd
}

func (cmd *VmSetCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}

	// if Help flag is set, no further processing is necessary
	if cmd.Help {
		return nil
	}

	// validate arg count
	cmdArgs, _ := cmd.Args()
	nExpected := 3
	if len(cmdArgs) != nExpected {
		cmd.PrintUsage()
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
			cmdString,
			nExpected,
			len(cmdArgs))
	}

	// init positional params (nop - no params)
	cmd.Name = cmdArgs[0]
	cmd.Key = cmdArgs[1]
	cmd.Value = cmdArgs[2]

	return nil
}

func (cmd *VmSetCommand) PrintUsage() {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Vm_Set)
	fmt.Println()
}

func (cmd VmSetCommand) Run() error {
	slog.Debug("VmSetCommand.Run()", "args", cmd.Cmdline)

	// parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}

	// If help flag set, print usage + return
	if cmd.Help {
		cmd.PrintUsage()
		return nil
	}

	// execute command
	err = RunVmSet(cmd.Name, cmd.Key, cmd.Value)
	return err
}

func RunVmSet(name string, key string, val string) error {
	slog.Debug("RunVmSet()", "name", name, "key", key, "val", val)
	// get vm-specific etcd client
	cli, err := eco.CreateVmClientFromConfig()
	if err != nil {
		return err
	}
	// get the vm data from the cluster, set the field (if it exists), and save updated object
	vm, err := cli.Get(name)
	if err != nil {
		return err
	}
	err = vm.Set(key, val)
	if err != nil {
		return err
	}
	vm, err = cli.Put(name, vm)
	if err != nil {
		return err
	}
	// print the updated vm if it exists
	if vm == nil {
		fmt.Printf("vm '%s' not found", name)
	} else {
		fmt.Println(vm)
	}

	return nil
}

// type VmShowCommand struct {
// 	*flag.FlagSet
// }

// func NewVmShowCommand () *VmShowCommand {
// 	// create command
// 	flagSet := flag.NewFlagSet("vm.show", flag.ExitOnError)
// 	cmd := VmShowCommand{FlagSet: flagSet}
// 	// init flag vars (nop -- no flags)

// 	return &cmd
// }

// func (cmd VmShowCommand) Run (args[] string) error {
// 	slog.Debug("VmShowCommand.Run()", "args", args)
// 	err := cmd.Parse(args)
// 	if err != nil { return err }
// 	return nil
// }

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
