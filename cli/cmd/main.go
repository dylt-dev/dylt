package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/dylt-dev/dylt/common"
)

type MainCommand struct {
	*BaseCommand
	Help bool // flag
}

func NewMainCommand(cmdline Cmdline) *MainCommand {
	flagSet := flag.NewFlagSet("dylt", flag.ExitOnError)
	var cmd = MainCommand{BaseCommand: &BaseCommand{Cmdline: cmdline, FlagSet: flagSet}}
	flagSet.BoolVar(&cmd.Help, "help", false, "give it to me")

	return &cmd
}

func (cmd *MainCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	return err
}

func (cmd *MainCommand) PrintUsage() {
	PrintUsage(USG_Main)
}

func (cmd *MainCommand) Run() error {
	if common.Logger == nil {
		panic("common.Logger == nil !!!")
	}
	common.Logger.Signature("MainCommand.Run()", cmd.Cmdline)
	// Check for 0 args; if so print usage & return
	// Parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}
	if cmd.Help {
		PrintUsage(USG_Main)
		return nil
	}
	if len(cmd.Cmdline) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}
	// Execute command
	err = RunMain(cmd.SubCommand(), cmd.SubArgs())
	if err != nil {
		return err
	}

	return nil
}

func RunMain(subCommand string, subCmdArgs []string) error {
	slog.Debug("RunMain()", "subCommand", subCommand, "subCmdArgs", subCmdArgs)

	// If there's no subcommand, do main() things
	if subCommand == "" {
		// Check if it's the user's first time. If so, act accordingly.
		is, err := isFirstTime()
		slog.Debug("main", "isFirstTime()", is)
		if err != nil {
			return err
		}
		if is {
			fmt.Println("Running firstTime() ...")
			err = firstTime()
			if err != nil {
				return err
			}
		}
	} else {
		// Create the subcommand and run it
		subCmd, err := createMainSubCommand(subCommand, subCmdArgs)
		if err != nil {
			return err
		}
		err = subCmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func createMainSubCommand(sCmd string, subCmdArgs []string) (Command, error) {
	switch sCmd {
	case "call":
		return NewCallCommand(subCmdArgs), nil
	case "config":
		return NewConfigCommand(subCmdArgs), nil
	case "get":
		return NewGetCommand(subCmdArgs), nil
	case "host":
		return NewHostCommand(subCmdArgs), nil
	case "init":
		return NewInitCommand(subCmdArgs), nil
	case "list":
		return NewListCommand(subCmdArgs), nil
	case "misc":
		return NewMiscCommand(subCmdArgs), nil
	case "status":
		return NewStatusCommand(subCmdArgs), nil
	case "vm":
		return NewVmCommand(subCmdArgs), nil
	case "watch":
		return NewWatchCommand(subCmdArgs), nil
	default:
		{
			var nilPtr *MainCommand = nil
			nilPtr.PrintUsage()
			return nil, fmt.Errorf("unrecognized subcommand: %s", sCmd)
		}
	}
}

func firstTime() error {
	fmt.Println("Welcome!")
	fmt.Println()
	fmt.Print("Please enter a domain for your etcd cluster -> ")
	// This is the user's first time.
	r := bufio.NewReader(os.Stdin)
	etcdDomain, err := r.ReadString('\n')
	etcdDomain = strings.TrimSpace(etcdDomain)
	if err != nil {
		return common.NewError(err)
	}
	cfg := common.ConfigStruct{EtcdDomain: etcdDomain}
	err = common.SaveConfig(cfg)
	if err != nil {
		return common.NewError(err)
	}
	fmt.Println("Saved!")
	fmt.Println()

	return nil
}

func isFirstTime() (bool, error) {
	configFilePath := common.GetConfigFilePath()
	common.Logger.Debugf("%s=%s", "configFilePath", configFilePath)
	_, err := os.Stat(configFilePath)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}

	if err == nil {
		common.Logger.Comment("config file found - not first time")
		return false, nil
	}

	// os.IsNotExist(err)
	common.Logger.Comment("config file does not exist - first time")
	return false, nil
}
