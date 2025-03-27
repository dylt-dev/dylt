package service

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dylt-dev/dylt/common"
	"github.com/dylt-dev/dylt/template"
)

func BuildDisableServiceCommand(svcName string) *exec.Cmd {
	return exec.Command("systemctl", "disable", svcName)
}
func BuildDoesServiceExistCommand(svcName string) *exec.Cmd {
	return exec.Command("systemctl", "is-enabled", svcName)
}
func BuildIsServiceActiveCommand(svcName string) *exec.Cmd {
	return exec.Command("systemctl", "is-active", svcName)
}
func BuildIsServiceEnabledCommand(svcName string) *exec.Cmd {
	return exec.Command("systemctl", "is-enabled", svcName)
}
func BuildStartServiceCommand(svcName string) *exec.Cmd {
	return exec.Command("systemctl", "start", svcName)
}
func BuildStopServiceCommand(svcName string) *exec.Cmd {
	return exec.Command("systemctl", "stop", svcName)
}

func BuildEnableServiceCommand(svcName string, svcFS *ServiceFS) *exec.Cmd {
	unitFilePath := svcFS.GetUnitFilePath(svcName)
	slog.Debug("BuildEnableServiceCommand", "unitFilePath", unitFilePath)
	cmd := exec.Command("systemctl", "enable", unitFilePath)
	return cmd
}

// simple typedef
type ServiceData map[string]string

// A service installed on a host.
// @note some ServiceFS methods take a service name. Seems like maybe
// if ServiceFS represents a service, then a svcName parameter doesn't
// make much sense.
type ServiceFS struct {
	Name     string
	RootPath string
}

func NewServiceFS(name string, rootPath string) *ServiceFS {
	return &ServiceFS{
		Name:     name,
		RootPath: rootPath,
	}
}

func (fs *ServiceFS) ChownSvc(uid int, gid int) error {
	folderPath := fs.GetFolderPath()
	err := common.ChownR(folderPath, uid, gid)
	if err != nil {
		return err
	}

	return nil
}

func (fs *ServiceFS) GetPath(filename string) string {
	return filepath.Join(fs.GetFolderPath(), filename)
}

func (fs *ServiceFS) GetFolderPath() string {
	return filepath.Join(fs.RootPath)
}

func (fs *ServiceFS) GetUnitFilePath(svcName string) string {
	filename := fmt.Sprintf("%s.service", svcName)
	return fs.GetPath(filename)
}

func (fs *ServiceFS) InitSvcFolder() error {
	// Create folder for service if necessary
	svcPath := fs.GetFolderPath()
	err := os.MkdirAll(svcPath, 0744)
	if err != nil {
		return err
	}

	return nil
}

func (fs *ServiceFS) OpenFile(filename string, flag int, perm os.FileMode) (*os.File, error) {
	path := fs.GetPath(filename)
	f, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (fs *ServiceFS) OpenWriter(filename string, perm os.FileMode) (*os.File, error) {
	return fs.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perm)
}

func (fs *ServiceFS) WriteRunScript(svc *ServiceSpec, tmpl *template.Template) error {
	path := tmpl.GetRunScriptPath()
	w, err := fs.OpenWriter(path, 0755)
	if err != nil {
		return err
	}
	defer w.Close()
	err = tmpl.WriteRunScript(w, svc.Data)
	if err != nil {
		return err
	}

	return nil
}

func (fs *ServiceFS) WriteUnitFile(svc *ServiceSpec, tmpl *template.Template) error {
	path := tmpl.GetUnitFilePath()
	w, err := fs.OpenWriter(path, 0644)
	if err != nil {
		return err
	}
	defer w.Close()
	err = tmpl.WriteUnitFile(w, svc.Data)
	if err != nil {
		return err
	}

	return nil
}

// A filesystem (fs.FS) that contains templates for each file of a service
// As such, it can generate the various files of a service
//
// @note now that I understand templates, template themselves are collections
// of templates so substructing Template might make make more sense than substructing
// FS.
//
// @note do I really only need one substruct for all services? Or will I end up with
// a different substruct for each type of service that's supported. Perhaps time
// will tell.
// type ServiceTemplateFS struct {
// 	fs.FS
// }

// func (fs *ServiceTemplateFS) GetRunScriptData(svcName string) (*TemplateData, error) {
// 	svcPattern := fmt.Sprintf("svc/%s/*", svcName)
// 	tmpl, err := template.ParseFS(fs.FS, svcPattern)
// 	if err != nil {
// 		return nil, err
// 	}
// 	runScriptFilename := "run.sh"
// 	tmpl = tmpl.Lookup(runScriptFilename)
// 	unitFileData := TemplateData{Template: *tmpl}

// 	return &unitFileData, nil
// }

// func (fs *ServiceTemplateFS) GetUnitFileData(svcName string) (*TemplateData, error) {
// 	svcPattern := fmt.Sprintf("svc/%s/*", svcName)
// 	tmpl, err := template.ParseFS(fs.FS, svcPattern)
// 	if err != nil {
// 		return nil, err
// 	}
// 	unitFilename := fmt.Sprintf("%s.service", svcName)
// 	tmpl = tmpl.Lookup(unitFilename)
// 	unitFileData := TemplateData{Template: *tmpl}

// 	return &unitFileData, nil
// }
