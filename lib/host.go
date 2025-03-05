package lib

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"text/template"
)

//go:embed svc/*
var FOL_Svc embed.FS
// const FN_WatchDaylightRunScript = "run.sh"
// const FN_WatchDaylightUnitFile = "watch-daylight.service"
// const PATH_WatchDaylightRunScript = "svc/watch-daylight/run.sh"
// const PATH_WatchDaylightUnitFile = "svc/watch-daylight/watch-daylight.service"
const UID_rayray = 2000
const GID_rayray = 2000
const PATH_SvcFolderRoot = "/opt/svc/"

func GetServiceFolder (svcName string) string {
	return filepath.Join(PATH_SvcFolderRoot, svcName)
}


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
			slog.Debug("lib.ChownR.fnWalk", "path", path, "d.Name()", d.Name())
			fullPath := filepath.Join(folderPath, path)
			err = os.Chown(fullPath, uid, gid)
			if err != nil { return err }
		}
		return nil
	}

	slog.Debug("lib.ChownR()", "folderPath", folderPath, "uid", uid, "gid", gid)
	dir := os.DirFS(folderPath)
	err := fs.WalkDir(dir, ".", fnWalk)
	if err != nil { return err }

	return nil
}

func ChownSvcFolder (svcName string, uid int, gid int) error {
	svcFolder := "/opt/svc"
	folder := filepath.Join(svcFolder, svcName)
	err := ChownR(folder, uid, gid)
	if err != nil { return err }

	return nil
}


func CreateWatchDaylightService (uid int, gid int) error {
	svcName := "watch-daylight"
	// Create folder for service if necessary
	err := InitSvcFolder(svcName)
	if err != nil { return err }
	// Open destination unit file & execute the template into the file
	data := ServiceData{}
	err = WriteUnitFile(svcName, data)
	if err != nil { return err }
	err = WriteRunScript(svcName, data)
	if err != nil { return err }
	// chown service folder to daylight user
	err = ChownSvcFolder("watch-daylight", uid, gid)
	if err != nil { return err }
	// run the service
	err = RunService(svcName)
	if err != nil { return err }
	
	return nil
}

func GetRunScriptData (svcName string) (*TemplateData, error) {
	svcPattern := fmt.Sprintf("svc/%s/*", svcName)
	tmpl, err := template.ParseFS(FOL_Svc, svcPattern)
	if err != nil { return nil, err }
	runScriptFilename := "run.sh"
	tmpl = tmpl.Lookup(runScriptFilename)
	unitFileData := TemplateData{Template: *tmpl}

	return &unitFileData, nil
}


func GetSvcWriter (svcName string, filename string, perm os.FileMode) (io.Writer, error) {
	svcFolder := GetServiceFolder(svcName)
	unitFilePath := path.Join(svcFolder, filename)
	w, err := os.OpenFile(unitFilePath, os.O_CREATE | os.O_WRONLY | os.O_TRUNC, perm)
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

func InitSvcFolder (svcName string) error {
	// Create folder for service if necessary
	svcFolder := GetServiceFolder(svcName)
	err := os.MkdirAll(svcFolder, 0744)
	if err != nil { return err }

	return nil
}

func RemoveService (svcName string) error {
	name := "systemctl"
	var cmd *exec.Cmd
	// systemctl stop $svcName
	slog.Debug("lib.RunService - exec.Command()", "args", []string{"stop", svcName})
	cmd = exec.Command(name, "stop", svcName)
	err := cmd.Run()
	if err != nil { return err }
	// systemctl disable $svcName
	slog.Debug("lib.RunService - execCommand", "args", []string{name, "disable", svcName})
	cmd = exec.Command(name, "disable", svcName)
	err = cmd.Run()
	if err != nil { return err }

	return nil
}

func RunService (svcName string) error {
	name := "systemctl"
	var cmd *exec.Cmd
	// systemctl daemon-reload
	slog.Debug("lib.RunService - exec.Command()", "args", []string{"daemon-reload"})
	cmd = exec.Command(name, "daemon-reload")
	err := cmd.Run()
	if err != nil { return err }
	// systemctl enable $unitFilePath
	unitFilename := fmt.Sprintf("%s.service", svcName)
	unitFilePath := filepath.Join(PATH_SvcFolderRoot, svcName, unitFilename)
	slog.Debug("lib.RunService - execCommand", "args", []string{name, "enable", unitFilePath})
	cmd = exec.Command(name, "enable", unitFilePath)
	err = cmd.Run()
	if err != nil { return err }
	// systemctl start $svcName
	slog.Debug("lib.RunService - exec.Command()", "args", []string{"start", svcName})
	cmd = exec.Command(name, "start", svcName)
	err = cmd.Run()
	if err != nil { return err }

	return nil
}

func WriteRunScript (svcName string, data ServiceData) error {
	runScriptData, err := GetRunScriptData(svcName)
	if err != nil { return err }
	runScriptFilename := "run.sh"
	w, err := GetSvcWriter(svcName, runScriptFilename, 0755)
	if err != nil { return err }
	err = runScriptData.Write(w, data)
	if err != nil { return err }

	return nil
}

func WriteUnitFile (svcName string, data ServiceData) error {
	unitFileData, err := GetUnitFileData(svcName)
	if err != nil { return err }
	unitFilename := fmt.Sprintf("%s.service", svcName)
	w, err := GetSvcWriter(svcName, unitFilename, 0644)
	if err != nil { return err }
	err = unitFileData.Write(w, data)
	if err != nil { return err }

	return nil
}