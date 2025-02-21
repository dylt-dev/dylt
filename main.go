package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path"

	// "github.com/dylt-dev/dylt/cli"
	// "github.com/dylt-dev/dylt/lib"
)

const LOG_File = "dylt.log"
const LOG_Folder = "/var/log/dylt"

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


func cmdInit () {
	fmt.Println("hello from cmdInit()")
}


func main () {
	// lib.InitConfig()
	initLogging()
	// os.Exit(cli.Run())
	// --token
	var token string
	flag.StringVar(&token, "token", "", "GitHub user access token")
	flag.Parse()
	args := flag.Args()
	slog.Debug("test", "args", args)
	cmd := args[0]
	switch cmd {
	case "init":
		cmdInit()
	}
}
