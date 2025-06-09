package cmd

import "fmt"

// dylt (main)
var USG_Main = fmt.Sprintf("\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n",
	USG_Call_Short,
	USG_Config_Short,
	USG_Get_Short,
	USG_Host_Short,
	USG_Init_Short,
	"list	(help)",
	"misc	(help)",
	"vm	(help)",
	"watch	(help)",
)


// dylt call
var USG_Call_Desc = "Invoke daylight.sh with command and optional args"
var USG_Call_Short = fmt.Sprintf("%-16s %s", "call", USG_Call_Desc)
var USG_Call = fmt.Sprintf("call [--script-path scriptPath] cmd [arg ... arg] %s", USG_Call_Desc)
var USG_Call_Full = fmt.Sprintf("\t%s\n\n\t%s\n",
	USG_Call,
	" --script-path (Optional) path to daylight.sh.",
)

// dylt config
const USG_Config_Desc = "get/set/show config values"
var USG_Config_Short = fmt.Sprintf("%-16s %s", "config", USG_Config_Desc)
var USG_Config = fmt.Sprintf("config subcmd [args], %s", USG_Config_Desc)
const USG_ConfigGet  = `config get key        get key from config`
const USG_ConfigSet  = `config set key val    set key to val in config`
const USG_ConfigShow = `config show           show current config contents`

// dylt get
const USG_Get_Desc = "get value from etcd"
var USG_Get_Short = createUsageShort("get", USG_Get_Desc)
var USG_Get = fmt.Sprintf("get key    %s", USG_Get_Desc)

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
var USG_Init_Desc = "initialize local daylight config data"
var USG_Init_Short = createUsageShort("init", USG_Init_Desc)
var USG_Init = createUsage("init --etcd-domain etcdDomain", USG_Init_Desc)

func createUsageShort (cmd string, desc string) string {
	return fmt.Sprintf("%-16s %s", cmd, desc)
}

func createUsage (cmdFull string, desc string) string {
	return fmt.Sprintf("%s    %s", cmdFull, desc)
}