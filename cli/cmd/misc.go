package cmd

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"os"
	"text/template"

	"github.com/dylt-dev/dylt/common"
	"github.com/dylt-dev/dylt/lib"
)

//go:embed content/*
var content embed.FS

type MiscCommand struct {
	*flag.FlagSet
	SubCommand string			// arg 0
	SubArgs    []string			// args 1..n-1
}

func NewMiscCommand () *MiscCommand {
	// create command
	flagSet := flag.NewFlagSet("misc", flag.ExitOnError)
	cmd := MiscCommand{FlagSet: flagSet}
	// init flag vars

	return &cmd
}

func (cmd *MiscCommand) HandleArgs(args []string) error {
	// parse flags
	err := cmd.Parse(args)
	if err != nil { return err }
	// validate arg count
	cmdArgs := cmd.Args()
	cmdName := "misc"
	nExpected := 1
	if len(cmdArgs) < nExpected {
		cmd.PrintUsage()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
			cmdName,
			nExpected,
			len(cmdArgs))
		}
	// init positional params
	cmd.SubCommand = cmdArgs[0]
	cmd.SubArgs = cmdArgs[1:]

	return nil
}

func (cmd *MiscCommand) PrintUsage () {
	fmt.Println()
	fmt.Printf("\t%s\n\t%s\n\t%s\n",
	USG_Misc_TwoNode_Short,
	USG_Misc_GenScript_Short,
	USG_Misc_Lookup_Short,
	)
	fmt.Println()
}

func (cmd *MiscCommand) Run(args []string) error {
	common.Logger.Debug("MiscCommand.Run()", "args", args)
	// parse flags & get positional args
	err := cmd.HandleArgs(args)
	if err != nil { return err }
	// execute command
	err = RunMisc(cmd.SubCommand, cmd.SubArgs)
	if err != nil { return err }

	return nil
}

func RunMisc (subCommand string, subCmdArgs []string) error {
	common.Logger.Debug("RunMisc()", "subCommand", subCommand, "subCmdArgs", subCmdArgs)
	// create the subcommand and run it
	subCmd, err := createMiscSubCommand(subCommand)
	if err != nil { return err }
	err = subCmd.Run(subCmdArgs)
	if err != nil { return err }

	return nil
}

func createMiscSubCommand (cmd string) (Command, error) {
	switch cmd {
	case "create-two-node-cluster": return NewCreateTwoNodeClusterCommand(), nil
	case "gen-etcd-run-script": return NewGenEtcdRunScriptCommand(), nil
	case "lookup": return NewLookupCommand(), nil
	default: return nil, fmt.Errorf("unrecognized subcommand: %s", cmd)
	}
}

type CreateTwoNodeClusterCommand struct {
	*flag.FlagSet
}

func NewCreateTwoNodeClusterCommand () *CreateTwoNodeClusterCommand {
	// create command
	flagSet := flag.NewFlagSet("createTwoNodeCluster", flag.ExitOnError)
	cmd := CreateTwoNodeClusterCommand{FlagSet: flagSet}
	// init flag vars

	return &cmd
}

func (cmd *CreateTwoNodeClusterCommand) HandleArgs(args []string) error {
	// parse flags
	err := cmd.Parse(args)
	if err != nil { return err }
	// validate arg count
	cmdArgs := cmd.Args()
	cmdName := "createTwoNodeCluster"
	nExpected := 0
	if len(cmdArgs) != nExpected {
		cmd.PrintUsage()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
			cmdName,
			nExpected,
			len(cmdArgs))
		}

	return nil
}

func (cmd *CreateTwoNodeClusterCommand) PrintUsage () {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Misc_TwoNode)
	fmt.Println()
}

func (cmd *CreateTwoNodeClusterCommand) Run(args []string) error {
	common.Logger.Debug("CreateTwoNodeClusterCommand.Run()", "args", args)
	// parse flags & get positional args
	err := cmd.HandleArgs(args)
	if err != nil { return err }
	// execute command
	// @getit
	err = RunCreateTwoNodeCluster()
	if err != nil { return err }

	return nil
}

func RunCreateTwoNodeCluster() error {
	common.Logger.Debug("RunCreateTwoNodeCluster()")
	var err error

	r := bufio.NewReader(os.Stdin)
	
	fmt.Println("Two node cluster time!")
	fmt.Println()

	fmt.Printf("Get your two node's IP addresses or hostnames, and whatever ssh private keys are necessary to connect to them. ")
	_, err = r.ReadBytes('\n')
	fmt.Println()

	fmt.Print("Done! (hit <Enter>) ")
	_, err = r.ReadBytes('\n')
	if err != nil { return err }

	return nil
}

type GenEtcdRunScriptCommand struct {
	*flag.FlagSet
}

func NewGenEtcdRunScriptCommand () *GenEtcdRunScriptCommand {
	// create command
	flagSet := flag.NewFlagSet("gen-etcd-run-script", flag.ExitOnError)
	cmd := GenEtcdRunScriptCommand{FlagSet: flagSet}
	// init flag vars

	return &cmd
}

func (cmd *GenEtcdRunScriptCommand) HandleArgs(args []string) error {
	// parse flags
	err := cmd.Parse(args)
	if err != nil { return err }
	// validate arg count
	cmdArgs := cmd.Args()
	cmdName := "gen-etcd-run-script"
	nExpected := 0
	if len(cmdArgs) != nExpected {
		cmd.PrintUsage()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
			cmdName,
			nExpected,
			len(cmdArgs))
		}

	return nil
}

func (cmd *GenEtcdRunScriptCommand) PrintUsage () {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Misc_GenScript)
	fmt.Println()
}

func (cmd *GenEtcdRunScriptCommand) Run(args []string) error {
	common.Logger.Debug("GenEtcdRunScriptCommand.Run()", "args", args)
	// parse flags & get positional args
	err := cmd.HandleArgs(args)
	if err != nil { return err }
	// execute command
	// @getit
	err = RunGenEtcdRunScript()
	if err != nil { return err }

	common.Logger.WithGroup("g")
	common.Logger.With("bar", "thirteen")
	common.Logger.Debug("testing logger", "foo", "13")
	return nil
}

func RunGenEtcdRunScript() error {
	common.Logger.Debug("RunGenEtcdRunScript()", )

	fmt.Println("I'm gennin a script!")

	buf, err := content.ReadFile("content/hello.tmpl")
	if err != nil { return err }
	tmpl := template.New("hello")
	tmpl, err = tmpl.Parse(string(buf))
	tmpl.Execute(os.Stdout, nil)
	return nil
}

type LookupCommand struct {
	*flag.FlagSet
	hostname string
}

func NewLookupCommand () *LookupCommand {
	flagSet := flag.NewFlagSet("misc lookup", flag.PanicOnError)
	cmd := LookupCommand { FlagSet: flagSet }

	return &cmd
}

func (cmd *LookupCommand) HandleArgs (args []string) error {
	// parse flags
	err := cmd.Parse(args)
	if err != nil { return err }
	// validate arg count
	cmdArgs := cmd.Args()
	cmdName := "lookup"
	nExpected := 1
	if len(cmdArgs) != nExpected {
		cmd.PrintUsage()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
			cmdName,
			nExpected,
			len(cmdArgs))
		}
	// init positional params
	cmd.hostname = cmdArgs[0]

	return nil
}

func (cmd *LookupCommand) PrintUsage () {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Misc_Lookup)
	fmt.Println()
}

func (cmd *LookupCommand) Run(args []string) error {
	common.Logger.Debug("LookupCommand.Run()", "args", args)
	// parse flags & get positional args
	err := cmd.HandleArgs(args)
	if err != nil { return err }
	// execute command
	// @getit
	err = lib.RunLookupCommand(cmd.hostname)
	if err != nil { return err }

	return nil
}
