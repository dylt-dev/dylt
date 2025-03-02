package cmd

import (
	"flag"
	"fmt"
	"log/slog"
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
	if len(cmdArgs) < nExpected { return fmt.Errorf("`%s` expects >=%d argument(s); received %d", cmdName, nExpected, len(cmdArgs)) }
	// init positional params
	cmd.SubCommand = cmdArgs[0]
	cmd.SubArgs = cmdArgs[1:]

	return nil
}

func (cmd *HostCommand) Run(args []string) error {
	slog.Debug("HostCommand.Run()", "args", args)
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
}

func NewHostInitCommand() *HostInitCommand {
	flagSet := flag.NewFlagSet("host.init", flag.ExitOnError)
	return &HostInitCommand{FlagSet: flagSet}
}

func (cmd *HostInitCommand) HandleArgs(args []string) error {
	// parse flags
	err := cmd.Parse(args)
	if err != nil { return err }
	// validate arg count
	cmdArgs := cmd.Args()
	cmdName := "host init"
	nExpected := 0
	if len(cmdArgs) != nExpected { return fmt.Errorf("`%s` expects %d argument(s); received %d", cmdName, nExpected, len(cmdArgs)) }
	// init positional params (nop - no positional params)

	return nil
}

func (cmd *HostInitCommand) Run(args []string) error {
	slog.Debug("HostInitCommand.Run()", "args", args)
	// Parse flags & get positional args
	err := cmd.HandleArgs(args)
	if err != nil { return err }
	// Execute command
	err = RunHostInit()
	if err != nil { return err }

	return nil
}

func RunHostInit() error {
	fmt.Println("init!")
	
	return nil
}
