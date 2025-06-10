package cmd

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/lib"
)

type CallCommand struct {
	*flag.FlagSet
	ScriptArgs []string		// args 0..n-1
	ScriptPath string		// flag
}

func NewCallCommand () *CallCommand {
	// create command
	flagSet := flag.NewFlagSet("call", flag.ExitOnError)
	cmd := CallCommand{FlagSet: flagSet}
	// init flag vars
	flagSet.StringVar(&cmd.ScriptPath, "script-path", "/opt/bin/daylight.sh", "script-path")
	
	return &cmd
}

func (cmd *CallCommand) HandleArgs (args []string) error {
	// parse flags
	err := cmd.Parse(args)
	if err != nil { return err }
	// validate arg count (nop - command takes all remaining args, 0 or more)
	cmdArgs := cmd.Args()
	Logger.Debugf("len(cmdArgs)=%d", len(cmdArgs))
	// Check for 0 args; if so print usage & return
	if len(cmdArgs) == 0 {
		Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// init positional params
	cmd.ScriptArgs = cmdArgs

	return nil
}

func (cmd *CallCommand) PrintUsage () {
	PrintMultilineUsage(USG_Call_Full)
}

func (cmd *CallCommand) Run (args []string) error {
	slog.Debug("CallCommand.Run()", "args", args)
	// Parse flags & get positional args
	err := cmd.HandleArgs(args)
	if err != nil { return err }
	// Execute command
	err = RunCall(cmd.ScriptPath, cmd.ScriptArgs)
	if err != nil { return err }

	return nil
}

func RunCall (scriptPath string, scriptArgs[] string) error {
	slog.Debug("RunCall()", "scriptPath", scriptPath, "scriptArgs", scriptArgs)
	// Call lib.RunScript() with script path and args, & output response
	_, s, err := lib.RunScript(scriptPath, scriptArgs)
	if err != nil { return err }
	fmt.Printf("%s\n", s)

	return nil
}