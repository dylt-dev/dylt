package cmd

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/dylt-dev/dylt/common"
)

type MainCommand struct {
	*BaseCommand
}

func NewMainCommand(cmdline Cmdline, parent Command) *MainCommand {
	// main command
	name := "dylt"
	cmdMap := CommandMap{
		"call": CallCommandF.New,
		"config": ConfigCommandF.New,
		"get": GetCommandF.New,
		"host": HostCommandF.New,
		"init": InitCommandF.New,
		"list": ListCommandF.New,
		"misc": MiscCommandF.New,
		"status": StatusCommandF.New,
		"vm": VmCommandF.New,
		"watch": WatchCommandF.New,
	}
	cmd := &MainCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Main, cmdMap)}

	//init flags (if any)

	return cmd
}

// func (cmd *MainCommand) CreateSubCommand() (Command, error) {
// 	args, flag := cmd.Args()
// 	if !flag {
// 		return nil, nil
// 	}
// 	return createMainSubCommand(args, cmd)
// }

func (cmd *MainCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	return err
}

func (cmd *MainCommand) Run() error {
	if common.Logger == nil {
		panic("common.Logger == nil !!!")
	}
	common.Logger.Signature("MainCommand.Run()", cmd.Cmdline)

	// Parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}

	// If help flag set, print usage
	if cmd.Help {
		cmd.PrintUsage()
		return nil
	}

	// If no args, print usage
	args, _ := cmd.Args()
	if len(args) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// Execute command
	err = RunMain(args, cmd)
	if err != nil {
		return err
	}

	return nil
}

func RunMain(cmdline Cmdline, cmd *MainCommand) error {
	slog.Debug("RunMain()", "cmdline", cmdline)

	// If there's no subcommand, do main() things
	if cmdline.Command() == "" {
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
		subCmd, err := cmd.CreateSubCommand()
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

// func createMainSubCommand(cmdline Cmdline, parent *MainCommand) (Command, error) {
// 	cmdName := cmdline.Command()
// 	cmdMap := CommandMap{
// 		"call": CallCommandF.New,
// 		"config": ConfigCommandF.New,
// 		"get": GetCommandF.New,
// 		"host": HostCommandF.New,
// 		"init": InitCommandF.New,
// 		"list": ListCommandF.New,
// 		"misc": MiscCommandF.New,
// 		"status": StatusCommandF.New,
// 		"vm": VmCommandF.New,
// 		"watch": WatchCommandF.New,
// 	}
	
// 	cmdFactoryFunc, ok := cmdMap[cmdName]
// 	if !ok {
// 		parent.PrintUsage()
// 		return nil, fmt.Errorf("unrecognized command: %s", cmdName)
// 	}
		
// 	cmd := cmdFactoryFunc(cmdline, parent)
// 	return cmd, nil
// }

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
