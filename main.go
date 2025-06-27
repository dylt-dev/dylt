package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	clicmd "github.com/dylt-dev/dylt/cli/cmd"
	"github.com/dylt-dev/dylt/common"
)

func exit (err error) {
	if err == nil {
		os.Exit(0)
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

func firstTime () error {
	fmt.Println("Welcome!")
	fmt.Println()
	fmt.Print("Please enter a domain for your etcd cluster -> ")
	// This is the user's first time.
	r := bufio.NewReader(os.Stdin)
	etcdDomain, err := r.ReadString('\n')
	etcdDomain = strings.TrimSpace(etcdDomain)
	if err != nil { return common.NewError(err) }
	cfg := common.ConfigStruct{ EtcdDomain: etcdDomain}
	err = common.SaveConfig(cfg)
	if err != nil { return common.NewError(err) }
	fmt.Println("Saved!")
	fmt.Println()

	return nil
}

func isFirstTime () (bool, error) {
	configFilePath := common.GetConfigFilePath()
	common.Logger.Debugf("%s=%s", "configFilePath", configFilePath)
	_, err := os.Stat(configFilePath)
	if err != nil && !os.IsNotExist(err) { return false, err }

	if err == nil {
		common.Logger.Comment("config file found - not first time")
		return false, nil
	}

	// os.IsNotExist(err)
	common.Logger.Comment("config file does not exist - first time")
	return false, nil
}

func createMainSubCommand (sCmd string) (clicmd.Command, error) {
	fmt.Printf("sCmd=%s\n", sCmd)
	switch sCmd {
	case "call": return clicmd.NewCallCommand(), nil
	case "config": return clicmd.NewConfigCommand(), nil
	case "get": return clicmd.NewGetCommand(), nil
	case "host": return clicmd.NewHostCommand(), nil
	case "init": return clicmd.NewInitCommand(), nil
	case "list": return clicmd.NewListCommand(), nil
	case "misc": return clicmd.NewMiscCommand(), nil
	case "status": return clicmd.NewStatusCommand(), nil
	case "vm": return clicmd.NewVmCommand(), nil
	case "watch": return clicmd.NewWatchCommand(), nil
	default: {
		var nilPtr *MainCommand = nil
		nilPtr.PrintUsage()
		return nil, fmt.Errorf("unrecognized subcommand: %s", sCmd)
	}
	}
}

type MainCommand struct {
	*flag.FlagSet
	SubCommand string
	SubArgs    []string
}

func NewMainCommand () *MainCommand {
	flagSet := flag.NewFlagSet("dylt", flag.ExitOnError)
	var cmd = MainCommand{FlagSet: flagSet}
	
	return &cmd
}

func (cmd *MainCommand) HandleArgs (args []string) error {
	// parse flags
	err := cmd.Parse(args)
	if err != nil { return err }
	// validate arg count - nop; command takes all remaining args, 0 or more)
	// init positional params - nop; there are none)

	return nil
}

func (cmd *MainCommand) PrintUsage () {
	clicmd.PrintMultilineUsage(clicmd.USG_Main)
}

func (cmd *MainCommand) Run (args clicmd.Cmdline) error {
	if common.Logger == nil { panic("common.Logger == nil !!!")}
	common.Logger.Signature("MainCommand.Run()", args)
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
	err = RunMain(args.Command(), args.Args())
	if err != nil { return err }

	return nil
}

func RunMain (subCommand string, subCmdArgs []string) error {
	slog.Debug("RunMain()", "subCommand", subCommand, "subCmdArgs", subCmdArgs)
	
	// If there's no subcommand, do main() things
	if subCommand == "" {
		// Check if it's the user's first time. If so, act accordingly.
		is, err := isFirstTime()
		slog.Debug("main", "isFirstTime()", is)
		if err != nil { exit(err) }
		if is {
			fmt.Println("Running firstTime() ...")
			err = firstTime()
			if err != nil { exit(err) }
		}
	} else {
		// Create the subcommand and run it
		subCmd, err := createMainSubCommand(subCommand)
		if err != nil { return err }
		err = subCmd.Run(subCmdArgs)
		if err != nil { return err }
	}

	return nil
}


func main () {
	common.InitLogging()
	
	cmd := NewMainCommand()
	var args clicmd.Cmdline = os.Args
	err := cmd.Run(args.Args())
	if err != nil {
		slog.Error(err.Error())
		fmt.Printf("\t%s\n", common.Error(err.Error()))
		fmt.Println()
		os.Exit(1)
	}
}
