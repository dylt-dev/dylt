package cmd

import (
	"errors"
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
	commandMap CommandMap
}

// type BaseCommandS BaseCommand[string]
// type BaseCommandSA BaseCommand[[]string]

func NewBaseCommand[U UsageTextType](name string, cmdline Cmdline, parent Command, usageText U, cmdMap CommandMap) *BaseCommand {
	cmd := &BaseCommand{
		Cmdline: cmdline,
		commandMap: cmdMap,
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

func (cmd BaseCommand) CommandLine() Cmdline {
	return cmd.Cmdline
}

func (cmd BaseCommand) CommandName() string {
	return cmd.Cmdline.Command()
}

func (cmd BaseCommand) CommandMap() CommandMap {
	fmt.Println("BaseCommand.CommandMap()")
	return cmd.commandMap
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
	fmt.Printf("%s %v\n", "cmd.CommandMap()", cmd.CommandMap())
	cmdline, is := cmd.Args()
	if !is {
		return nil, errors.New("Command not Parse()'d")
	}
	if cmd.CommandMap() == nil {
		return nil, nil
	}
	// return createConfigSubCommand(args, cmd)
	cmdName := cmdline.Command()
	cmdMap := cmd.CommandMap()
	if cmdMap == nil {
		return nil, nil
	}

	cmdFactoryFunc, ok := cmdMap[cmdName]
	if !ok {
		cmd.PrintUsage()
		return nil, fmt.Errorf("unrecognized command: %s", cmdName)
	}
		
	subCmd := cmdFactoryFunc(cmdline, cmd)
	return subCmd, nil
}

func (cmd BaseCommand) HandleArgs () error { return nil }


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

func (cmd BaseCommand) Run() error { return nil }

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
