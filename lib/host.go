package lib

import (
	"embed"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

//go:embed svc/*
var FOL_Svc embed.FS

const UID_rayray = 2000
const GID_rayray = 2000
const PATH_SvcFolderRoot = "/opt/svc/"

type TemplateData struct {
	template.Template
}

func (unitFile *TemplateData) Write(w io.Writer, data map[string]string) error {
	err := unitFile.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

func ChownR(folderPath string, uid int, gid int) error {
	var fnWalk fs.WalkDirFunc = func(path string, d fs.DirEntry, err error) error {
		if err == nil {
			slog.Debug("lib.ChownR.fnWalk", "path", path, "d.Name()", d.Name())
			fullPath := filepath.Join(folderPath, path)
			err = os.Chown(fullPath, uid, gid)
			if err != nil {
				return err
			}
		}
		return nil
	}

	slog.Debug("lib.ChownR()", "folderPath", folderPath, "uid", uid, "gid", gid)
	dir := os.DirFS(folderPath)
	err := fs.WalkDir(dir, ".", fnWalk)
	if err != nil {
		return err
	}

	return nil
}

func CreateWatchDaylightService(uid int, gid int) error {
	slog.Debug("lib.CreateWatchDaylightService()", "uid", uid, "gid", gid)
	const svcName = "watch-daylight"
	var svc ServiceSpec = ServiceSpec{svcName, ServiceData{}}
	var svcFS ServiceFS = ServiceFS{RootPath: PATH_SvcFolderRoot}
	var templateFS ServiceTemplateFS = ServiceTemplateFS{FS: FOL_Svc}
	var err error

	// Remove the service if it exists
	slog.Info("Removing service ...")
	err = RemoveService(svcName, &svcFS)
	if err != nil { return err }

	// install the service
	slog.Info("Installing service ...")
	err = InstallService(&svc, &templateFS, &svcFS, uid, gid)
	if err != nil { return err }

	// run the service
	slog.Info("Running service ...")
	err = RunService(svcName, &svcFS)
	if err != nil { return err }

	return nil
}

func CreateWatchSvcService(uid int, gid int) error {
	slog.Debug("lib.CreateWatchSvcService()", "uid", uid, "gid", gid)
	const svcName = "watch-svc"
	var svc ServiceSpec = ServiceSpec{svcName, ServiceData{}}
	var svcFS ServiceFS = ServiceFS{RootPath: PATH_SvcFolderRoot}
	var templateFS ServiceTemplateFS = ServiceTemplateFS{FS: FOL_Svc}
	var err error

	// Remove the service if it exists
	slog.Info("Removing service ...")
	err = RemoveService(svcName, &svcFS)
	if err != nil { return err }

	// install the service
	slog.Info("Installing service ...")
	err = InstallService(&svc, &templateFS, &svcFS, uid, gid)
	if err != nil { return err }

	// run the service
	slog.Info("Running service ...")
	err = RunService(svcName, &svcFS)
	if err != nil { return err }

	return nil
}

func InstallService (svc *ServiceSpec, templateFS *ServiceTemplateFS, svcFS *ServiceFS, uid int, gid int) error {
	var err error

	// Create folder for service if necessary
	slog.Info("Initializing service folder ...")
	err = svcFS.InitSvcFolder(svc.Name)
	if err != nil { return err }

	// Execute unit file template & write to file
	slog.Info("Writing Unit file ...")
	err = svcFS.WriteUnitFile(svc, templateFS)
	if err != nil { return err }

	// Execute run script template & write to file
	slog.Info("Writing run script ...")
	err = svcFS.WriteRunScript(svc, templateFS)
	if err != nil { return err }

	// chown service folder to daylight user
	slog.Info("Chown'ing service ...")
	err = svcFS.ChownSvc("watch-daylight", uid, gid)
	if err != nil { return err }

	return nil
}

func RemoveService(svcName string, fs *ServiceFS) error {
	slog.Debug("lib.RemoveService()", "svcName", svcName)
	var svc ServiceSpec = ServiceSpec{Name: svcName}
	var err error

	// Stop service
	slog.Info("Stopping service ...")
	err = svc.Stop()
	if err != nil { return err }

	// Disable service
	slog.Info("Disabling service ...")
	err = svc.Disable()
	if err != nil { return err }
	
	// Remove service folder
	slog.Info("Removing service ...")
	err = svc.Remove(fs)
	if err != nil { return err }

	return nil
}

func RunService(svcName string, svcFS *ServiceFS) error {
	slog.Debug("lib.RunService()", "svcName", svcName)
	var svc ServiceSpec = ServiceSpec{Name: svcName}
	var err error

	// systemctl daemon-reload
	slog.Info("Running `systemctl daemon-reload` ...")
	cmd := exec.Command("systemctl", "daemon-reload")
	err = cmd.Run()
	if err != nil { return err }

	// systemctl enable $unitFilePath
	slog.Info("Enabling service ...")
	err = svc.Enable(svcFS)
	if err != nil { return err }

	// systemctl start $svcName
	slog.Info("Starting service ...")
	err = svc.Start()
	if err != nil { return err }

	return nil
}