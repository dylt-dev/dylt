package lib

import (
	"embed"
	"fmt"
	"io"
)

//go:embed svc/*
var FOL_Svc embed.FS
const FN_WatchDaylightRunScript = "run.sh"
const FN_WatchDaylightUnitFile = "watch-daylight.service"
const PATH_WatchDaylightRunScript = "svc/watch-daylight/run.sh"
const PATH_WatchDaylightUnitFile = "svc/watch-daylight/watch-daylight.service"


func CreateWatchDaylightService () error {
	f, err := FOL_Svc.Open(PATH_WatchDaylightUnitFile)
	if err != nil { return err }
	s, err := io.ReadAll(f)
	if err != nil { return err }
	fmt.Println(string(s))

	return nil
}