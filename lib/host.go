package lib

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"path"
)

//go:embed svc/*
var FOL_Svc embed.FS
const FN_WatchDaylightRunScript = "run.sh"
const FN_WatchDaylightUnitFile = "watch-daylight.service"
const PATH_WatchDaylightRunScript = "svc/watch-daylight/run.sh"
const PATH_WatchDaylightUnitFile = "svc/watch-daylight/watch-daylight.service"
const UID_rayray = 2000
const GID_rayray = 2000


func CreateWatchDaylightService () error {
	// Create folder for service if necessary
	svcFolder := "/opt/svc/watch-daylight-go"
	err := os.MkdirAll(svcFolder, 0744)
	if err != nil { panic(fmt.Sprintf("Couldn't create or open service folder: %s", svcFolder)) }
	// Get unit file template from embedded FS
	tmpl, err := template.ParseFS(FOL_Svc, "svc/watch-daylight/*")
	if err != nil { return err }
	tmpl = tmpl.Lookup(FN_WatchDaylightUnitFile)
	unitFilename := "watch-daylight.service"
	unitFilePath := path.Join(svcFolder, unitFilename)
	// Open destination unit file & execute the template into the file
	unitFile, err := os.OpenFile(unitFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil { return err }
	data := map[string]string{}
	err = tmpl.Execute(unitFile, data)
	if err != nil { return err }
	// chown service folder to rayray

	return nil
}