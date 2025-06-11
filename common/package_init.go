package common

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"reflect"

	"github.com/dylt-dev/dylt/color"
)

const LOG_File = "dylt.log"
const LOG_Folder = "/var/log/dylt"

var Logger *cliLogger

func init () {
	// Logger = NewLogger(os.Stdout)
}

func FullTypeName(ty reflect.Type) string {
	var typeName = ty.Name()
	if typeName == "" {
		typeName = "(anon)"
	}

	pkgPath := ty.PkgPath()
	if pkgPath == "" {
		return typeName
	}

	if filepath.Dir(pkgPath) == "github.com/dylt-dev/dylt" {
		pkgPath = filepath.Base(pkgPath)
	}

	return fmt.Sprintf("%s.%s", pkgPath, typeName)
}

func Error (s string) string {

	var ss = color.Styledstring(s).Fg(color.Sys.Red)

	return string(ss)
}

func Highlight(s string) string {

	var ss = color.Styledstring(s).Style("\033[1;97m")

	return string(ss)
}

func InitLogging () {
	logToFile, ok := os.LookupEnv("DYLT_LogToFile")

	if ok  && logToFile != "" {
		InitLoggingToFile()
	} else {
		var logger = NewLogger(os.Stdout)
		setLogger(logger)
	}
	slog.Debug("logging successfully initialized")
}


func InitLoggingToFile () {
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

	// opts := slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}
	// var logger *slog.Logger = slog.New(slog.NewJSONHandler(logWriter, &opts))
	var logger = NewLogger(logWriter)
	setLogger(logger)

	// slog.SetDefault(logger)
	slog.Debug("logging successfully initialized")
}


func Lowlight(s string) string {
	var ss = color.Styledstring(s).Fg(color.X11.Gray50)

	return string(ss)
}

// All data specified in flags to the `init` subcommand
type InitStruct struct {
	EtcdDomain string
}

func Init(initData *InitStruct) error {
	cfg := ConfigStruct{
		EtcdDomain: initData.EtcdDomain,
	}
	err := SaveConfig(cfg)
	if err != nil { return err }

	return nil
}

func setLogger (logger *cliLogger) {
	Logger = logger
	slog.SetDefault(Logger.Logger)
}