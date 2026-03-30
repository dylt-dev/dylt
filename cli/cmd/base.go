package cmd

import (
	"flag"
	"fmt"
	"strings"
)

type BaseCommand struct {
	*flag.FlagSet
	Parent  Command
	Cmdline Cmdline
	Help    bool
	Usage string
}

// type BaseCommandS BaseCommand[string]
// type BaseCommandSA BaseCommand[[]string]

func NewBaseCommand[U UsageTextType](name string, cmdline Cmdline, parent Command, usageText U) *BaseCommand {
	cmd := &BaseCommand{
		Cmdline: cmdline,
		Parent:  parent,
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),
		Usage: CreateUsageString(usageText),
	}
	cmd.FlagSet.BoolVar(&cmd.Help, "help", false, "give it to me")

	return cmd
}

func (cmd BaseCommand) Args() (Cmdline, bool) {
	if !cmd.FlagSet.Parsed() {
		return nil, false
	}
	return cmd.FlagSet.Args(), true
}

func (cmd BaseCommand) CommandLine() Cmdline {
	return cmd.Cmdline
}

func (cmd BaseCommand) CommandName() string {
	return cmd.Cmdline.Command()
}

func (cmd BaseCommand) CommandArgs() ([]string, bool) {
	if !cmd.FlagSet.Parsed() {
		return nil, false
	}
	var cmdArgs []string
	var flag bool
	if cmd.Parent == nil {
		cmdArgs = []string{}
	} else {
		cmdArgs, flag = cmd.Parent.CommandArgs()
		if !flag {
			return nil, flag
		}
	}
	cmdArgs = append(cmdArgs, cmd.Cmdline.Command())
	return cmdArgs, true
}

func (cmd BaseCommand) CommandString() (string, bool) {
	if !cmd.FlagSet.Parsed() {
		return "", false
	}
	cmdArgs, flag := cmd.CommandArgs()
	if !flag {
		return "", false
	}
	if cmdArgs == nil {
		return "", flag
	}
	return strings.Join(cmdArgs, " "), true
}

type NoSubcommandsError struct {}
func (o NoSubcommandsError) Error() string { return "No Subcommands" }

func (cmd BaseCommand) CreateSubCommand () (Command, error) {
	return nil, &NoSubcommandsError{}
}

func (cmd BaseCommand) Parse() error {
	err := cmd.FlagSet.Parse(cmd.Cmdline.Args())
	if err != nil {
		return err
	}
	return nil
}

func (cmd BaseCommand) PrintUsage () {
	fmt.Print(cmd.Usage)
	fmt.Println()
}

func (cmd BaseCommand) SubArgs() (Cmdline, bool) {
	if !cmd.FlagSet.Parsed() {
		return nil, false
	}
	var subCmdline Cmdline = cmd.FlagSet.Args()
	return subCmdline.Args(), true
}

func (cmd BaseCommand) SubCommand() (string, bool) {
	if !cmd.Parsed() {
		return "", false
	}
	var subCmdline Cmdline = cmd.FlagSet.Args()
	return subCmdline.Command(), true
}

type ICommandFactory[U Command] interface {
	New (Cmdline, Command) U
}

type CommandFactory[U Command] struct {
	FnNew func (Cmdline, Command) U
}

func (cf CommandFactory[U]) New (cmdline Cmdline, parent Command) U {
	return cf.FnNew(cmdline, parent)
}

