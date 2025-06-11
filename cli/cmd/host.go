package cmd

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/api"
	"github.com/dylt-dev/dylt/common"
)

type HostCommand struct {
	*flag.FlagSet
	SubCommand string
	SubArgs    []string

}

func NewHostCommand () *HostCommand {
	// create command
	flagSet := flag.NewFlagSet("host", flag.ExitOnError)
	cmd := HostCommand{FlagSet: flagSet}
	// init flag vars (nop -- no flags)

	return &cmd
}

func (cmd *HostCommand) HandleArgs(args []string) error {
	// parse flags
	err := cmd.Parse(args)
	if err != nil { return err }
	// validate arg count
	cmdArgs := cmd.Args()
	cmdName := "host"
	nExpected := 1
	if len(cmdArgs) < nExpected {
		cmd.PrintUsage()
		return fmt.Errorf("`%s` expects >=%d argument(s); received %d", cmdName, nExpected, len(cmdArgs))
	}
	// init positional params
	cmd.SubCommand = cmdArgs[0]
	cmd.SubArgs = cmdArgs[1:]

	return nil
}

func (cmd *HostCommand) PrintUsage () {
	PrintMultilineUsage(USG_Host)
}

func (cmd *HostCommand) Run(args []string) error {
	slog.Debug("HostCommand.Run()", "args", args)
	// Check for 0 args; if so print usage & return
	if len(args) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}
	// Parse flags & get positional args
	err := cmd.HandleArgs(args)
	if err != nil { return err }
	// Execute command
	err = RunHost(cmd.SubCommand, cmd.SubArgs)
	if err != nil { return err }

	return nil
}

func RunHost(subCommand string, subCmdArgs []string) error {
	slog.Debug("RunHost()", "subCommand", subCommand, "subCmdArgs", subCmdArgs)
	// Create the subcommand and run it
	subCmd, err := createHostSubCommand(subCommand)
	if err != nil { return err }
	err = subCmd.Run(subCmdArgs)
	if err != nil { return err }

	return nil
}

func createHostSubCommand(cmdName string) (Command, error) {
	switch cmdName {
		case "init": return NewHostInitCommand(), nil
	default:
		return nil, fmt.Errorf("unrecognized command: %s", cmdName)
	}
}

type HostInitCommand struct {
	*flag.FlagSet
	Gid int
	Uid int
}

func NewHostInitCommand() *HostInitCommand {
	flagSet := flag.NewFlagSet("host.init", flag.ExitOnError)
	cmd := HostInitCommand{FlagSet: flagSet}
	flagSet.IntVar(&cmd.Gid, "gid", 2000, "gid")
	flagSet.IntVar(&cmd.Uid, "uid", 2000, "uid")

	return &cmd
}

func (cmd *HostInitCommand) HandleArgs(args []string) error {
	// parse flags
	err := cmd.Parse(args)
	if err != nil { return err }
	// validate arg count
	cmdArgs := cmd.Args()
	cmdName := "host init"
	nExpected := 0
	if len(cmdArgs) != nExpected {
		cmd.PrintUsage()
		return fmt.Errorf("`%s` expects %d argument(s); received %d", cmdName, nExpected, len(cmdArgs))
	}
	// init positional params (nop - no positional params)

	return nil
}

func (cmd *HostInitCommand) PrintUsage () {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Host_Init)
	fmt.Println()
}

func (cmd *HostInitCommand) Run(args []string) error {
	slog.Debug("HostInitCommand.Run()", "args", args)
	// Parse flags & get positional args
	err := cmd.HandleArgs(args)
	if err != nil { return err }
	// Execute command
	err = RunHostInit(cmd.Uid, cmd.Gid)
	if err != nil { return err }

	return nil
}

func RunHostInit(uid int, gid int) error {
	slog.Debug("RunHostInit()", "uid", uid, "gid", gid)
	var err error
	fmt.Println("init!")

	err = api.CreateWatchDaylightService(uid, gid)
	if err != nil { return err }
	
	err = api.CreateWatchSvcService(uid, gid)
	if err != nil { return err }

	return nil
}
