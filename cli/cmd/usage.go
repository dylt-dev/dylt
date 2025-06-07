package cmd

import "fmt"

// dylt (main)
var USG_Main = fmt.Sprintf("\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n",
	USG_Call_Short,
	USG_Config_Short,
	"get	(help)",
	"host	(help)",
	"init	(help)",
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
var USG_Config_Desc = "get/set/show config values"
var USG_Config_Short = fmt.Sprintf("%-16s %s", "config", USG_Config_Desc)
var USG_Config = fmt.Sprintf("config subcmd [args], %s", USG_Config_Desc)
const USG_ConfigGet  = `config get key        get key from config`
const USG_ConfigSet  = `config set key val    set key to val in config`
const USG_ConfigShow = `config show           show current config contents`

