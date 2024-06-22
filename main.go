package main

import (
	"log/slog"
	"os"
	"path"

	"github.com/dylt-dev/dylt/cli"
	"github.com/dylt-dev/dylt/lib"
)


func initLogging () {
	logFile := "dylt.log"
	logFolder := "/var/log/dylt"
	logPath := path.Join(logFolder, logFile)
	logWriter, err := os.OpenFile(logPath, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0777)
	if err != nil { panic("Couldn't open logfile") }
	opts := slog.HandlerOptions{AddSource: true}
	var logger *slog.Logger = slog.New(slog.NewJSONHandler(logWriter, &opts))
	slog.SetDefault(logger)
}


func main () {
	lib.InitConfig()
	initLogging()
	os.Exit(cli.Run())
}
