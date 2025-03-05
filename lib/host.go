package lib

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

//go:embed svc/*
var FOL_Svc embed.FS
const FN_WatchDaylightRunScript = "run.sh"
const FN_WatchDaylightUnitFile = "watch-daylight.service"
const PATH_WatchDaylightRunScript = "svc/watch-daylight/run.sh"
const PATH_WatchDaylightUnitFile = "svc/watch-daylight/watch-daylight.service"
const UID_rayray = 2000
const GID_rayray = 2000


type ServiceData map[string]string
type TemplateData struct {
	template.Template
}


func (unitFile *TemplateData) Write (w io.Writer, data map[string]string) error {
	err := unitFile.Execute(w, data)
	if err != nil { return err }

	return nil
}

func ChownR (folderPath string, uid int, gid int) error {
	var fnWalk fs.WalkDirFunc = func (path string, d fs.DirEntry, err error) error {
		if err == nil {
			slog.Debug("lib.ChownR", "path", path, "d.Name()", d.Name())
			fullPath := filepath.Join(folderPath, path)
			err = os.Chown(fullPath, uid, gid)
			if err != nil { return err }
		}
		return nil
	}

	dir := os.DirFS(folderPath)
	err := fs.WalkDir(dir, ".", fnWalk)
	if err != nil { return err }

	return nil
}

func ChownSvcFolder (svcName string, uid int, gid int) error {
	svcFolder := "/opt/svc/watch-daylight-go"
	folder := filepath.Join(svcFolder, svcName)
	err := ChownR(folder, uid, gid)
	if err != nil { return err }

	return nil
}


func CreateWatchDaylightService (uid int, gid int) error {
	svcName := "watch-daylight"
	// Create folder for service if necessary
	err := InitSvcFolder()
	if err != nil { return err }
	// Open destination unit file & execute the template into the file
	data := map[string]string{}
	err = WriteUnitFile(svcName, data)
	if err != nil { return err }
	// chown service folder to daylight user
	err = ChownSvcFolder("watch-daylight", uid, gid)
	if err != nil { return err }

	return nil
}

func GetRunScriptData (svcName string) (*TemplateData, error) {
	svcPattern := fmt.Sprintf("svc/%s/*", svcName)
	tmpl, err := template.ParseFS(FOL_Svc, svcPattern)
	if err != nil { return nil, err }
	runScriptFilename := fmt.Sprintf("%s.service", svcName)
	tmpl = tmpl.Lookup(runScriptFilename)
	unitFileData := TemplateData{Template: *tmpl}

	return &unitFileData, nil
}


func GetSvcWriter (filename string) (io.Writer, error) {
	svcFolder := "/opt/svc/watch-daylight-go"
	unitFilePath := path.Join(svcFolder, filename)
	w, err := os.OpenFile(unitFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil { return nil, err }

	return w, nil
}

func GetUnitFileData (svcName string) (*TemplateData, error) {
	svcPattern := fmt.Sprintf("svc/%s/*", svcName)
	tmpl, err := template.ParseFS(FOL_Svc, svcPattern)
	if err != nil { return nil, err }
	unitFilename := fmt.Sprintf("%s.service", svcName)
	tmpl = tmpl.Lookup(unitFilename)
	unitFileData := TemplateData{Template: *tmpl}

	return &unitFileData, nil
}

func InitSvcFolder () error {
	// Create folder for service if necessary
	svcFolder := "/opt/svc/watch-daylight-go"
	err := os.MkdirAll(svcFolder, 0744)
	if err != nil { return err }

	return nil
}

func WriteRunScript (svcName string, data ServiceData) error {
	runScriptData, err := GetRunScriptData(svcName)
	if err != nil { return err }
	runScriptFilename := "run.sh"
	w, err := GetSvcWriter(runScriptFilename)
	if err != nil { return err }
	err = runScriptData.Write(w, data)
	if err != nil { return err }

	return nil
}

func WriteUnitFile (svcName string, data ServiceData) error {
	unitFileData, err := GetUnitFileData(svcName)
	if err != nil { return err }
	unitFilename := fmt.Sprintf("%s.service", svcName)
	w, err := GetSvcWriter(unitFilename)
	if err != nil { return err }
	err = unitFileData.Write(w, data)
	if err != nil { return err }

	return nil
}