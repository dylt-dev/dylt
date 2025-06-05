package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strings"

	clicmd "github.com/dylt-dev/dylt/cli/cmd"
	"github.com/dylt-dev/dylt/common"
)

const LOG_File = "dylt.log"
const LOG_Folder = "/var/log/dylt"


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

func initLogging () {
	logToFile, ok := os.LookupEnv("DYLT_LogToFile")
	if ok  && logToFile != "" {
		initLoggingToFile()
	} else {
		opts := slog.HandlerOptions{AddSource: false, Level: slog.LevelDebug}
		var logger *slog.Logger = slog.New(slog.NewTextHandler(os.Stdout, &opts))
		slog.SetDefault(logger)
		slog.Debug("logging successfully initialized")
	}
}


func initLoggingToFile () {
	var logFile, logFolder, logPath string
	envLogPath, ok := os.LookupEnv("DYLT_LogPath")
	if ok {
		logFile = path.Base(envLogPath)
		logFolder = path.Dir(envLogPath)
		logPath = envLogPath
	} else {
		logFile = LOG_File
		logFolder = LOG_Folder
		logPath = path.Join(logFolder, logFile)
	}
	err := os.MkdirAll(logFolder, 0744)
	if err != nil { panic(fmt.Sprintf("Couldn't create or open log folder: %s", logFolder)) }
	logWriter, err := os.OpenFile(logPath, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0777)
	if err != nil { panic(fmt.Sprintf("Couldn't open logfile path: logFile=%s logFolder=%s logPath=%s", logFile, logFolder, logPath)) }
	opts := slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}
	var logger *slog.Logger = slog.New(slog.NewJSONHandler(logWriter, &opts))
	slog.SetDefault(logger)
	slog.Debug("logging successfully initialized")
}


func isFirstTime () (bool, error) {
	configFilePath := common.GetConfigFilePath()
	_, err := os.Stat(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func createMainSubCommand (cmd string) (clicmd.Command, error) {
	switch cmd {
	case "call": return clicmd.NewCallCommand(), nil
	case "config": return clicmd.NewConfigCommand(), nil
	case "get": return clicmd.NewGetCommand(), nil
	case "host": return clicmd.NewHostCommand(), nil
	case "init": return clicmd.NewInitCommand(), nil
	case "list": return clicmd.NewListCommand(), nil
	case "misc": return clicmd.NewMiscCommand(), nil
	case "vm": return clicmd.NewVmCommand(), nil
	case "watch": return clicmd.NewWatchCommand(), nil
	default: return nil, fmt.Errorf("unrecognized subcommand: %s", cmd)
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

func (cmd *MainCommand) Run (args clicmd.Cmdline) error {
	clicmd.Logger.Debug("CallCommand.Run()", "args", args)
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
	
	initLogging()
	
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
	cmd := NewMainCommand()
	var args clicmd.Cmdline = os.Args
	err := cmd.Run(args.Args())
	if err != nil {
		slog.Error(err.Error())
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
