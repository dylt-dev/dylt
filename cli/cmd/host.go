package cmd

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/api"
	"github.com/dylt-dev/dylt/common"
)

type HostCommand struct {
	*BaseCommand
}

func NewHostCommand(cmdline Cmdline, parent Command) *HostCommand {
	// create command
	flagSet := flag.NewFlagSet("host", flag.ExitOnError)
	cmd := HostCommand{BaseCommand: &BaseCommand{Cmdline: cmdline, FlagSet: flagSet, ParentCommand: parent}}
	// init flag vars (nop -- no flags)

	return &cmd
}

func (cmd *HostCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}
	// validate arg count
	cmdArgs, _ := cmd.Args()
	nExpected := 1
	if len(cmdArgs) < nExpected {
		cmd.PrintUsage()
		cmdString, _ := cmd.GetCommandString()
		return fmt.Errorf("`%s` expects >=%d argument(s); received %d",
		                  cmdString,
						  nExpected,
						  len(cmdArgs))
	}
	// init positional params (nop - no positional params)

	return nil
}

func (cmd *HostCommand) PrintUsage() {
	PrintUsage(USG_Host)
}

func (cmd *HostCommand) Run() error {
	slog.Debug("HostCommand.Run()", "args", cmd.Cmdline)
	// Check for 0 args; if so print usage & return
	if len(cmd.Cmdline) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}
	// Parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}
	// Execute command
	subArgs, _ := cmd.SubArgs()
	subCommand, _ := cmd.SubCommand()
	err = RunHost(subCommand, subArgs)
	if err != nil {
		return err
	}

	return nil
}

func RunHost(subCommand string, subCmdArgs []string) error {
	slog.Debug("RunHost()", "subCommand", subCommand, "subCmdArgs", subCmdArgs)
	// Create the subcommand and run it
	subCmd, err := createHostSubCommand(subCommand, subCmdArgs)
	if err != nil {
		return err
	}
	err = subCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func createHostSubCommand(cmdName string, subCmdArgs Cmdline) (Command, error) {
	switch cmdName {
	case "init":
		return NewHostInitCommand(subCmdArgs), nil
	default:
		return nil, fmt.Errorf("unrecognized command: %s", cmdName)
	}
}

type HostInitCommand struct {
	*BaseCommand
	Gid int
	Uid int
}

func NewHostInitCommand(cmdline Cmdline) *HostInitCommand {
	flagSet := flag.NewFlagSet("host.init", flag.ExitOnError)
	cmd := HostInitCommand{BaseCommand: &BaseCommand{FlagSet: flagSet}}
	flagSet.IntVar(&cmd.Gid, "gid", 2000, "gid")
	flagSet.IntVar(&cmd.Uid, "uid", 2000, "uid")

	return &cmd
}

func (cmd *HostInitCommand) HandleArgs() error {
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
		cmdString, _ := cmd.GetCommandString()
		return fmt.Errorf("%s` expects %d argument(s); received %d",
		                  cmdString,
						  nExpected,
						  len(cmdArgs))
	}
	// init positional params (nop - no positional params)

	return nil
}

func (cmd *HostInitCommand) PrintUsage() {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Host_Init)
	fmt.Println()
}

func (cmd *HostInitCommand) Run() error {
	slog.Debug("HostInitCommand.Run()", "args", cmd.Cmdline)
	// Parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}
	// Execute command
	err = RunHostInit(cmd.Uid, cmd.Gid)
	if err != nil {
		return err
	}

	return nil
}

func RunHostInit(uid int, gid int) error {
	slog.Debug("RunHostInit()", "uid", uid, "gid", gid)
	var err error
	fmt.Println("init!")

	err = api.CreateWatchDaylightService(uid, gid)
	if err != nil {
		return err
	}

	err = api.CreateWatchSvcService(uid, gid)
	if err != nil {
		return err
	}

	return nil
}
