package cmd

import "fmt"

// dylt (main)
var USG_Main = fmt.Sprintf("\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n",
	USG_Call_Short,
	USG_Config_Short,
	USG_Get_Short,
	USG_Host_Short,
	USG_Init_Short,
	USG_List_Short,
	USG_Misc_Short,
	USG_Vm_Short,
	"watch	(help)",
)


// dylt call
var USG_Call_Desc = "Invoke daylight.sh with command and optional args"
var USG_Call_Short = createUsageShort("call", USG_Call_Desc)
var USG_Call = fmt.Sprintf("call [--script-path scriptPath] cmd [arg ... arg] %s", USG_Call_Desc)
var USG_Call_Full = fmt.Sprintf("\t%s\n\n\t%s\n",
	USG_Call,
	" --script-path (Optional) path to daylight.sh.",
)

// dylt config
const USG_Config_Desc = "get/set/show config values"
var USG_Config_Short = createUsageShort("config", USG_Config_Desc)
var USG_Config = fmt.Sprintf("config subcmd [args], %s", USG_Config_Desc)
const USG_ConfigGet  = `config get key        get key from config`
const USG_ConfigSet  = `config set key val    set key to val in config`
const USG_ConfigShow = `config show           show current config contents`

// dylt get
const USG_Get_Desc = "get value from etcd"
var USG_Get_Short = createUsageShort("get", USG_Get_Desc)
var USG_Get = createUsage("get key", USG_Get_Desc)

// dylt host
// var USG_Host_Desc = "host host hosty host host"
// var USG_Host_Short = createUsageShort("host", USG_Host_Desc)
var USG_Host_Desc = "host stuff"
var USG_Host_Short = createUsageShort("host", USG_Host_Desc)
var USG_Host = fmt.Sprintf("\t%s\n",
	USG_Host_Init_Short,
)

// dylt host init
const USG_Host_Init_Desc = "prepare a host for daylight"
var USG_Host_Init_Short = createUsageShort("host init", USG_Host_Init_Desc)
var USG_Host_Init = createUsage("host init", USG_Host_Init_Desc)

// dylt init
const USG_Init_Desc = "initialize local daylight config data"
var USG_Init_Short = createUsageShort("init", USG_Init_Desc)
var USG_Init = createUsage("init --etcd-domain etcdDomain", USG_Init_Desc)

// dylt list
const USG_List_Desc = "list all keys in cluster"
var USG_List_Short = createUsageShort("list", USG_List_Desc)
var USG_List = createUsage("list", USG_List_Desc)

// dylt misc
const USG_Misc_Desc = "Miscellaneous collection of commands"
var USG_Misc_Short = createUsageShort("misc", USG_Misc_Desc)

// dylt misc create-two-node-cluster
const USG_Misc_TwoNode_Desc = "Create a simple 2 node daylight cluster"
var USG_Misc_TwoNode_Short = createUsageShort("misc create-two-node-cluster", USG_Misc_TwoNode_Desc)

// dylt misc gen-etcd-run-script
const USG_Misc_GenScript_Desc = "Generate a script for running etcd"
var USG_Misc_GenScript_Short = createUsageShort("misc gen-etcd-run-script", USG_Misc_GenScript_Desc)

// dylt vm
const USB_Vm_Desc = "VM-related operations"
var USG_Vm_Short = createUsageShort("vm", USB_Vm_Desc)

// dy;t vm add
const USG_Vm_Add_Desc = "create a new VM"
var USG_Vm_Add_Short = createUsageShort("vm add", USG_Vm_Add_Desc)
var USG_Vm_Add = createUsage("vm add name fqdn", USG_Vm_Add_Desc)

// dylt vm all
const USG_Vm_All_Desc = "list all VMs"
var USG_Vm_All_Short = createUsageShort("vm all", USG_Vm_All_Desc)
var USG_Vm_All = createUsage("vm all", USG_Vm_All_Desc)

// dylt vm del
const USG_Vm_Del_Desc = "delete a VM"
var USG_Vm_Del_Short = createUsageShort("vm del", USG_Vm_Del_Desc)
var USG_Vm_Del = createUsage("vm del vmName", USG_Vm_Del_Desc)

// dylt vm get
const USG_Vm_Get_Desc = "get data for a VM"
var USG_Vm_Get_Short = createUsageShort("vm get", USG_Vm_Get_Desc)
var USG_Vm_Get = createUsage("vm get vmName", USG_Vm_Get_Desc)

// dylt vm list
const USG_Vm_List_Desc = "list all VMs"
var USG_Vm_List_Short = createUsageShort("vm list", USG_Vm_List_Desc)
var USG_Vm_List = createUsage("vm list", USG_List_Desc)

// dylt vm set
const USG_Vm_Set_Desc = "set values for a VM"
var USG_Vm_Set_Short = createUsageShort("vm set", USG_Vm_Set_Desc)
var USG_Vm_Set = createUsage("vm set vmName key val", USG_Vm_Set_Desc)

func createUsageShort (cmd string, desc string) string {
	return fmt.Sprintf("%-16s # %s", cmd, desc)
}

func createUsage (cmdFull string, desc string) string {
	return fmt.Sprintf("%s    # %s", cmdFull, desc)
}