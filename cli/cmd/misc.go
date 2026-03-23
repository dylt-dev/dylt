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
	*BaseCommand
}

func NewMiscCommand(cmdline Cmdline, parent SuperCommand) *MiscCommand {
	// create command
	flagSet := flag.NewFlagSet("misc", flag.ExitOnError)
	cmd := MiscCommand{BaseCommand: &BaseCommand{Cmdline: cmdline, FlagSet: flagSet, ParentCommand: parent}}
	// init flag vars - no flags; nop

	return &cmd
}

func (cmd *MiscCommand) CreateSubCommand () (Command, error) {
	args, is := cmd.Args()
	if !is {
		return nil, nil
	}
	return createMiscSubCommand(args, cmd)
}

func (cmd *MiscCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}
	// validate arg count
	nExpected := 0
	if len(cmd.Cmdline) < nExpected {
		cmd.PrintUsage()
		return fmt.Errorf("`%s` expects >=%d argument(s); received %d",
			cmd.Cmdline[0],
			nExpected,
			len(cmd.Cmdline))
	}
	// init positional params

	return nil
}

func (cmd *MiscCommand) PrintUsage() {
	fmt.Println()
	fmt.Printf("\t%s\n\t%s\n\t%s\n",
		USG_Misc_TwoNode_Short,
		USG_Misc_GenScript_Short,
		USG_Misc_Lookup_Short,
	)
	fmt.Println()
}

func (cmd *MiscCommand) Run() error {
	common.Logger.Debug("MiscCommand.Run()", "args", cmd.Cmdline)
	// parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}
	// If no args, print usage
	args, _ := cmd.Args()
	if len(args) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}
	// execute command
	err = RunMisc(args, cmd)
	return err
}

func RunMisc(cmdline Cmdline, parent SuperCommand) error {
	common.Logger.Debug("RunMisc()", "cmdline", cmdline, "parent", parent)
	// create the subcommand and run it
	subCmd, err := createMiscSubCommand(cmdline, parent)
	if err != nil {
		return err
	}
	err = subCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func createMiscSubCommand(cmdline Cmdline, parent SuperCommand) (Command, error) {
	cmdName := cmdline.Command()
	switch cmdName {
	case "create-two-node-cluster":
		return NewCreateTwoNodeClusterCommand(cmdline, parent), nil
	case "gen-etcd-run-script":
		return NewGenEtcdRunScriptCommand(cmdline, parent), nil
	case "lookup":
		return NewLookupCommand(cmdline, parent), nil
	default:
		return nil, fmt.Errorf("unrecognized subcommand: %s", cmdName)
	}
}

type CreateTwoNodeClusterCommand struct {
	*BaseCommand
}

func NewCreateTwoNodeClusterCommand(cmdline Cmdline, parent SuperCommand) *CreateTwoNodeClusterCommand {
	// create command
	flagSet := flag.NewFlagSet("createTwoNodeCluster", flag.ExitOnError)
	cmd := CreateTwoNodeClusterCommand{BaseCommand: &BaseCommand{Cmdline: cmdline, FlagSet: flagSet, ParentCommand: parent}}
	// init flag vars

	return &cmd
}

func (cmd *CreateTwoNodeClusterCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}
	// validate arg count
	cmdArgs, _ := cmd.Args()
	nExpected := 0
	if len(cmd.Cmdline) < nExpected {
		cmd.PrintUsage()
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
			cmdString,
			nExpected,
			len(cmdArgs))
	}

	return nil
}

func (cmd *CreateTwoNodeClusterCommand) PrintUsage() {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Misc_TwoNode)
	fmt.Println()
}

func (cmd *CreateTwoNodeClusterCommand) Run() error {
	common.Logger.Debug("CreateTwoNodeClusterCommand.Run()", "args", cmd.Cmdline)
	// parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}
	// execute command
	// @getit
	err = RunCreateTwoNodeCluster()
	return err
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
	return err
}

type GenEtcdRunScriptCommand struct {
	*BaseCommand
}

func NewGenEtcdRunScriptCommand(cmdline Cmdline, parent SuperCommand) *GenEtcdRunScriptCommand {
	// create command
	flagSet := flag.NewFlagSet("gen-etcd-run-script", flag.ExitOnError)
	cmd := GenEtcdRunScriptCommand{BaseCommand: &BaseCommand{Cmdline: cmdline, FlagSet: flagSet, ParentCommand: parent}}
	// init flag vars - no flags; nop

	return &cmd
}

func (cmd *GenEtcdRunScriptCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
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

	return nil
}

func (cmd *GenEtcdRunScriptCommand) PrintUsage() {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Misc_GenScript)
	fmt.Println()
}

func (cmd *GenEtcdRunScriptCommand) Run() error {
	common.Logger.Debug("GenEtcdRunScriptCommand.Run()", "args", cmd.Cmdline)
	// parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}
	// execute command
	// @getit
	err = RunGenEtcdRunScript()
	if err != nil {
		return err
	}

	common.Logger.WithGroup("g")
	common.Logger.With("bar", "thirteen")
	common.Logger.Debug("testing logger", "foo", "13")
	return nil
}

func RunGenEtcdRunScript() error {
	common.Logger.Debug("RunGenEtcdRunScript()")

	fmt.Println("I'm gennin a script!")

	buf, err := content.ReadFile("content/hello.tmpl")
	if err != nil {
		return err
	}
	tmpl := template.New("hello")
	tmpl, err = tmpl.Parse(string(buf))
	tmpl.Execute(os.Stdout, nil)
	return nil
}

type LookupCommand struct {
	*BaseCommand
	Hostname string
}

func NewLookupCommand(cmdline Cmdline, parent SuperCommand) *LookupCommand {
	flagSet := flag.NewFlagSet("misc lookup", flag.PanicOnError)
	cmd := LookupCommand{BaseCommand: &BaseCommand{Cmdline: cmdline, FlagSet: flagSet, ParentCommand: parent}}

	return &cmd
}

func (cmd *LookupCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
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
	cmd.Hostname = cmdArgs[0]

	return nil
}

func (cmd *LookupCommand) PrintUsage() {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Misc_Lookup)
	fmt.Println()
}

func (cmd *LookupCommand) Run() error {
	common.Logger.Debug("LookupCommand.Run()", "args", cmd.Cmdline)
	// parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}
	// If no args, print usage
	if len(cmd.Cmdline) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}
	// execute command
	// @getit
	err = lib.RunLookupCommand(cmd.Hostname)
	return err
}
