package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	clicmd "github.com/dylt-dev/dylt/cli/cmd"
	"github.com/dylt-dev/dylt/common"
)

func main() {
	common.InitLogging()

	var args clicmd.Cmdline = os.Args
	cmd := clicmd.NewMainCommand(args, nil)
	err := cmd.Run()
	if err != nil {
		// common.PrintBlankIfTerminal(os.Stderr)
		// // yikes - need a cleaner idiom for logging errors
		// common.Logger.Error(common.Error(err.Error()))
		// common.PrintBlankIfTerminal(os.Stderr)
		exit(err)
	}
}

func exit(err error) {
	if err == nil {
		// os.Exit(0)
	}
	slog.Error(err.Error())
	fmt.Println(err.Error())
	switch err := err.(type) {
	case *exec.ExitError:
		os.Exit(err.ExitCode())
	default:
		os.Exit(1)
	}
}
