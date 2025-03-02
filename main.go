package main

import (
	"embed"
	"fmt"
	"log/slog"
	"os"
	"path"

	clicmd "github.com/dylt-dev/dylt/cli/cmd"
)

const LOG_File = "dylt.log"
const LOG_Folder = "/var/log/dylt"

//go:embed svc
var svc embed.FS

func initLogging () {
	err := os.MkdirAll(LOG_Folder, 0744)
	if err != nil { panic("Couldn't create or open log folder") }
	logPath := path.Join(LOG_Folder, LOG_File)
	logWriter, err := os.OpenFile(logPath, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0777)
	if err != nil { panic("Couldn't open logfile") }
	opts := slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}
	var logger *slog.Logger = slog.New(slog.NewJSONHandler(logWriter, &opts))
	slog.SetDefault(logger)
}


func createSubCommand (cmd string) (clicmd.Command, error) {
	switch cmd {
	case "call": return clicmd.NewCallCommand(), nil
	case "config": return clicmd.NewConfigCommand(), nil
	case "get": return clicmd.NewGetCommand(), nil
	case "init": return clicmd.NewInitCommand(), nil
	case "list": return clicmd.NewListCommand(), nil
	case "vm": return clicmd.NewVmCommand(), nil
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
