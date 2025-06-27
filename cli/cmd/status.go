package cmd

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/lib"
)

type StatusCommand struct {
	*flag.FlagSet
}

func NewStatusCommand () *StatusCommand {
	// create command
	flagSet := flag.NewFlagSet("status", flag.ExitOnError)
	cmd := StatusCommand{FlagSet: flagSet}
	// init flag vars

	return &cmd
}

func (cmd *StatusCommand) HandleArgs(args []string) error {
	// parse flags
	err := cmd.Parse(args)
	if err != nil { return err }
	// validate arg count
	cmdArgs := cmd.Args()
	cmdName := "status"
	nExpected := 0
	if len(cmdArgs) != nExpected {
		return fmt.Errorf("`%s` expects %d argument(s); received %d",
			cmdName,
			nExpected,
			len(cmdArgs))
		}

	return nil
}

func (cmd *StatusCommand) Run(args []string) error {
	slog.Debug("StatusCommand.Run()", "args", args)
	// parse flags & get positional args
	err := cmd.HandleArgs(args)
	if err != nil { return err }
	// execute command
	// @getit
	err = RunStatus()
	if err != nil { return err }

	return nil
}

func RunStatus() error {
	slog.Debug("RunStatus()")

	lib.RunStatus()

	return nil
}