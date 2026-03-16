package cmd

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/common"
	"github.com/dylt-dev/dylt/lib"
)

type CallCommand struct {
	*BaseCommand
	ScriptArgs []string // args 0..n-1
	ScriptPath string   // flag
}

func NewCallCommand(cmdline Cmdline) *CallCommand {
	// create command
	flagSet := flag.NewFlagSet("call", flag.ExitOnError)
	cmd := CallCommand{BaseCommand: &BaseCommand{Cmdline: cmdline, FlagSet: flagSet}}
	// init flag vars
	flagSet.StringVar(&cmd.ScriptPath, "script-path", "/opt/bin/daylight.sh", "script-path")

	return &cmd
}

func (cmd *CallCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}
	// validate arg count (nop - command takes all remaining args, 0 or more)
	common.Logger.Debugf("len(cmd.Cmdline)=%d", len(cmd.Cmdline))
	// Check for 0 args; if so print usage & return
	if len(cmd.Cmdline) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// init positional params
	cmd.ScriptArgs = cmd.Cmdline

	return nil
}

func (cmd *CallCommand) PrintUsage() {
	PrintUsage(USG_Call_Full)
}

func (cmd *CallCommand) Run() error {
	slog.Debug("CallCommand.Run()", "args", cmd.SubArgs())
	// Parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}
	// Execute command
	err = RunCall(cmd.ScriptPath, cmd.ScriptArgs)
	if err != nil {
		return err
	}

	return nil
}

func RunCall(scriptPath string, scriptArgs []string) error {
	slog.Debug("RunCall()", "scriptPath", scriptPath, "scriptArgs", scriptArgs)
	// Call lib.RunScript() with script path and args, & output response
	_, s, err := lib.RunScript(scriptPath, scriptArgs)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", s)

	return nil
}
