package cmd

import (
	"flag"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
)

type BaseCommand struct {
	*flag.FlagSet
	ParentCommand Command
	Cmdline       Cmdline
	Help          bool
}

func NewBaseCommand (name string, cmdline Cmdline, parent SuperCommand) *BaseCommand {
	cmd := &BaseCommand{
		Cmdline: cmdline,
		ParentCommand: parent,
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),
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

func (cmd BaseCommand) CommandLine () Cmdline {
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
	if cmd.ParentCommand == nil {
		cmdArgs = []string{}
	} else {
		cmdArgs, flag = cmd.ParentCommand.CommandArgs()
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

func (cmd BaseCommand) Log() {
	stype := reflect.TypeOf(cmd).Name()
	subArgs, _ := cmd.SubArgs()
	slog.Debug(fmt.Sprintf("%s.Run()", stype), "args", subArgs)
}

func (cmd BaseCommand) Parse() error {
	err := cmd.FlagSet.Parse(cmd.Cmdline.Args())
	if err != nil {
		return err
	}
	return nil
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
