package cmd

import (
	"flag"
)

type BaseCommand struct {
	*flag.FlagSet
	Cmdline Cmdline
}


func (cmd *BaseCommand) Parse() error {
	err := cmd.FlagSet.Parse(cmd.Cmdline)
	if err != nil {
		return err
	}
	return nil
}

func (cmd *BaseCommand) SubArgs() []string {
	if !cmd.Parsed() {
		return []string{}
	}
	var subCmdline Cmdline = cmd.Args()
	return subCmdline.Args()
}

func (cmd *BaseCommand) SubCommand() string {
	if !cmd.Parsed() {
		return ""
	}
	var subCmdline Cmdline = cmd.Args()
	return subCmdline.Command()
}


