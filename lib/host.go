package lib

import (
	"embed"
	"log/slog"
	"os/exec"

	"github.com/dylt-dev/dylt/service"
	"github.com/dylt-dev/dylt/template"
)

//go:embed svcfiles/*
var EMBED_SvcFiles embed.FS

const DEF_uid_rayray = 2000
const DEF_gid_rayray = 2000
const DEF_SvcFolderRootPath = "/opt/svc/"


func CreateWatchDaylightService(uid int, gid int) error {
	slog.Debug("lib.CreateWatchDaylightService()", "uid", uid, "gid", gid)
	const svcName = "watch-daylight"
	var svc service.ServiceSpec = service.ServiceSpec{Name: svcName, Data: service.ServiceData{}}
	var svcFS service.ServiceFS = service.ServiceFS{RootPath: DEF_SvcFolderRootPath}
	var tmpl *template.Template = template.New(svcName)
	var err error

	// remove the service if it exists
	slog.Info("Removing service ...")
	err = RemoveService(svcName, &svcFS)
	if err != nil {
		return err
	}

	// install the service
	slog.Info("Installing service ...")
	err = InstallService(&svc, tmpl, &svcFS, uid, gid)
	if err != nil {
		return err
	}

	// run the service
	slog.Info("Running service ...")
	err = RunService(svcName, &svcFS)
	if err != nil {
		return err
	}

	return nil
}

func CreateWatchSvcService(uid int, gid int) error {
	slog.Debug("lib.CreateWatchSvcService()", "uid", uid, "gid", gid)
	const svcName = "watch-svc"
	var svc service.ServiceSpec = service.ServiceSpec{Name: svcName, Data: service.ServiceData{}}
	var svcFS *service.ServiceFS = service.NewServiceFS(svcName, DEF_SvcFolderRootPath)
	var tmpl *template.Template = template.New(svcName)
	var err error

	// Remove the service if it exists
	slog.Info("Removing service ...")
	err = RemoveService(svcName, svcFS)
	if err != nil {
		return err
	}

	// install the service
	slog.Info("Installing service ...")
	err = InstallService(&svc, tmpl, svcFS, uid, gid)
	if err != nil {
		return err
	}

	// run the service
	slog.Info("Running service ...")
	err = RunService(svcName, svcFS)
	if err != nil {
		return err
	}

	return nil
}

func InstallService(svc *service.ServiceSpec, tmpl *template.Template, svcFS *service.ServiceFS, uid int, gid int) error {
	var err error

	// Create folder for service if necessary
	slog.Info("Initializing service folder ...")
	err = svcFS.InitSvcFolder()
	if err != nil {
		return err
	}

	// Execute unit file template & write to file
	slog.Info("Writing Unit file ...")
	err = svcFS.WriteUnitFile(svc, tmpl)
	if err != nil {
		return err
	}

	// Execute run script template & write to file
	slog.Info("Writing run script ...")
	err = svcFS.WriteRunScript(svc, tmpl)
	if err != nil {
		return err
	}

	// chown service folder to daylight user
	slog.Info("Chown'ing service ...")
	err = svcFS.ChownSvc(uid, gid)
	if err != nil {
		return err
	}

	return nil
}

func RemoveService(svcName string, fs *service.ServiceFS) error {
	slog.Debug("lib.RemoveService()", "svcName", svcName)
	var svc service.ServiceSpec = service.ServiceSpec{Name: svcName}
	var err error

	// Stop service
	slog.Info("Stopping service ...")
	err = svc.Stop()
	if err != nil {
		return err
	}

	// Disable service
	slog.Info("Disabling service ...")
	err = svc.Disable()
	if err != nil {
		return err
	}

	// Remove service folder
	slog.Info("Removing service ...")
	err = svc.Remove(fs)
	if err != nil {
		return err
	}

	return nil
}

func RunService(svcName string, svcFS *service.ServiceFS) error {
	slog.Debug("lib.RunService()", "svcName", svcName)
	var svc *service.ServiceSpec = service.NewServiceSpec(svcName)
	var err error

	// systemctl daemon-reload
	slog.Info("Running `systemctl daemon-reload` ...")
	cmd := exec.Command("systemctl", "daemon-reload")
	err = cmd.Run()
	if err != nil {
		return err
	}

	// systemctl enable $unitFilePath
	slog.Info("Enabling service ...")
	err = svc.Enable(svcFS)
	if err != nil {
		return err
	}

	// systemctl start $svcName
	slog.Info("Starting service ...")
	err = svc.Start()
	if err != nil {
		return err
	}

	return nil
}
