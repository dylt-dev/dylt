package main

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	clicmd "github.com/dylt-dev/dylt/cli/cmd"
)

const LOG_File = "dylt.log"
const LOG_Folder = "/var/log/dylt"

func initLogging () {
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


func createSubCommand (cmd string) (clicmd.Command, error) {
	switch cmd {
	case "call": return clicmd.NewCallCommand(), nil
	case "config": return clicmd.NewConfigCommand(), nil
	case "get": return clicmd.NewGetCommand(), nil
	case "host": return clicmd.NewHostCommand(), nil
	case "init": return clicmd.NewInitCommand(), nil
	case "list": return clicmd.NewListCommand(), nil
	case "vm": return clicmd.NewVmCommand(), nil
	case "watch": return clicmd.NewWatchCommand(), nil
	default: return nil, fmt.Errorf("unrecognized subcommand: %s", cmd)
	}
}

func main () {
	// lib.InitConfig()
	initLogging()
	var cmdline clicmd.Cmdline = os.Args[1:]
	cmdName := cmdline.Command()
	slog.Debug("main", "cmdName", cmdName)
	cmdArgs := cmdline.Args()
	cmd, err := createSubCommand(cmdName)
	if err != nil { os.Exit(1) }
	err = cmd.Run(cmdArgs)
	if err != nil {
		slog.Error(err.Error())
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
