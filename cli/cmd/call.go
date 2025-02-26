package cmd

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/lib"
)

type CallCommand struct {
	*flag.FlagSet
	ScriptArgs []string
	ScriptPath string
}

func NewCallCommand () *CallCommand {
	flagSet := flag.NewFlagSet("call", flag.ExitOnError)
	cmd := CallCommand{FlagSet: flagSet}
	flagSet.StringVar(&cmd.ScriptPath, "script-path", "/opt/bin/daylight.sh", "script-path")
	return &cmd
}

func (cmd *CallCommand) Run (args []string) error {
	// Parse flags & get positional args
	err := cmd.Parse(args)
	if err != nil { return err }
	cmd.ScriptArgs = cmd.Args()
	// Execute command
	return RunCall(cmd.ScriptPath, cmd.ScriptArgs)
}

func RunCall (scriptPath string, scriptArgs[] string) error {
	slog.Debug("In RunCall()", "scriptPath", scriptPath, "scriptArgs", scriptArgs)
	_, stdout, err := lib.RunScript(scriptPath, scriptArgs)
	if err != nil { return err }
	fmt.Printf("%s\n", stdout)
	return nil
}