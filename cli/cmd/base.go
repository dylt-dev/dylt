package cmd

import (
	"flag"
	"fmt"
	"log/slog"
	"reflect"
)

type BaseCommand struct {
	*flag.FlagSet
	ParentCommand Command
	Cmdline Cmdline
}

func (cmd *BaseCommand) Args () (Cmdline, bool) {
	if !cmd.FlagSet.Parsed() {
		return nil, false
	}
	return cmd.FlagSet.Args(), true
}

func (cmd *BaseCommand) GetCommandString () string {
	return ""
}

func (cmd *BaseCommand) Log() {
	stype := reflect.TypeOf(cmd).Name()
	subArgs, _ := cmd.SubArgs()
	slog.Debug(fmt.Sprintf("%s.Run()", stype), "args", subArgs)
}

func (cmd *BaseCommand) Parse() error {
	err := cmd.FlagSet.Parse(cmd.Cmdline)
	if err != nil {
		return err
	}
	return nil
}

func (cmd *BaseCommand) SubArgs() (Cmdline, bool) {
	if !cmd.FlagSet.Parsed() {
		return nil, false
	}
	var subCmdline Cmdline = cmd.FlagSet.Args()
	return subCmdline.Args().Args(), true
}

func (cmd *BaseCommand) SubCommand() (string, bool) {
	if !cmd.Parsed() {
		return "", false
	}
	var subCmdline Cmdline = cmd.FlagSet.Args()
	return subCmdline.Args().Command(), true
}
