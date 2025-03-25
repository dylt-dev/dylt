package lib

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
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
	Name string
	RootPath string
}

func NewServiceFS (name string, rootPath string) *ServiceFS {
	return &ServiceFS{
		Name: name,
		RootPath: rootPath,
	}
}

func (fs* ServiceFS) ChownSvc(uid int, gid int) error {
	folderPath := fs.GetFolderPath()
	err := ChownR(folderPath, uid, gid)
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

func (fs* ServiceFS) WriteRunScript(svc *ServiceSpec, templateFS *ServiceTemplateFS) error {
	runScriptData, err := templateFS.GetRunScriptData(svc.Name)
	if err != nil {
		return err
	}
	runScriptFilename := "run.sh"
	w, err := fs.OpenWriter(runScriptFilename, 0755)
	if err != nil {
		return err
	}
	defer w.Close()
	err = runScriptData.Write(w, svc.Data)
	if err != nil {
		return err
	}

	return nil
}

func (fs* ServiceFS) WriteUnitFile(svc *ServiceSpec, templateFS *ServiceTemplateFS) error {
	unitFileData, err := templateFS.GetUnitFileData(svc.Name)
	if err != nil {
		return err
	}
	unitFilename := fmt.Sprintf("%s.service", svc.Name)
	w, err := fs.OpenWriter(unitFilename, 0644)
	if err != nil {
		return err
	}
	defer w.Close()
	err = unitFileData.Write(w, svc.Data)
	if err != nil {
		return err
	}

	return nil
}

// Struct representing a systemd service. It consists of a service's name,
// and the data making up the service.
// @note this might be a tad redundant with ServiceData, since it's just ServiceData
// plus a name. ServiceData is a map, giving this the illusion this makes sense
type ServiceSpec struct {
	Name string
	Data ServiceData
}

func NewServiceSpec(name string) *ServiceSpec {
	spec := ServiceSpec{
		Name: name,
		Data: ServiceData{},
	}

	return &spec
}

// eg systemctl disable $service
func (svc *ServiceSpec) BuildDisableCommand () *exec.Cmd {
	return BuildDisableServiceCommand(svc.Name)
}

// eg systemctl is-enabled $service
func (svc *ServiceSpec) BuildDoesExistCommand () *exec.Cmd {
	return BuildDoesServiceExistCommand(svc.Name)
}

// eg systemctl enable /path/to/unit/file.service
func (svc *ServiceSpec) BuildEnableCommand (svcFS ServiceFS) *exec.Cmd {
	return BuildEnableServiceCommand(svc.Name, &svcFS)
}

// eg systemctl is-active $service
func (svc *ServiceSpec) BuildIsActiveCommand() *exec.Cmd {
	return BuildIsServiceActiveCommand(svc.Name)
}

// eg systemctl is-enabled $service
func (svc* ServiceSpec) BuildIsEnabledCommand() *exec.Cmd {
	return BuildIsServiceEnabledCommand(svc.Name)
}

// eg systemctl start $service
func (svc *ServiceSpec) BuildStartCommand () *exec.Cmd {
	return BuildStartServiceCommand(svc.Name)
}

// eg systemctl stop $service
func (svc *ServiceSpec) BuildStopCommand () *exec.Cmd {
	return BuildStopServiceCommand(svc.Name)
}

func (svc *ServiceSpec) Disable () error {
	// Check if service exists. If not, leave peacefully.
	does, err := svc.IsExists()
	if err != nil { return err }
	if !does { slog.Debug(fmt.Sprintf("%s doesn't exist.", svc.Name)); return nil }

	// Stop service if running
	err = svc.Stop()
	if err != nil { return err }

	// Check if service is enabled - if not, leave peacefully
	is, err := svc.IsEnabled()
	slog.Debug(fmt.Sprintf("is=%t err=%s", is, err))
	if err != nil { return err }
	if !is { fmt.Printf("%s is already disabled", svc.Name); return nil }

	// Build and run command
	slog.Debug(fmt.Sprintf("Disabling %s ...", svc.Name))
	cmd := svc.BuildDisableCommand()
	err = cmd.Run()
	if err != nil { return err }

	return nil
}

func (svc *ServiceSpec) Enable (svcFS *ServiceFS) error {
	// Check if service is enabled. If so, leave peacefully.
	is, err := svc.IsEnabled()
	if err != nil { return err }
	if is { slog.Debug(fmt.Sprintf("%s is already enabled", svc.Name)); return nil }
	
	// Build and run command
	slog.Debug(fmt.Sprintf("Enable %s ...", svc.Name))
	cmd := svc.BuildEnableCommand(*svcFS)
	err = cmd.Run()
	if err != nil { return err }

	return nil
}

// Answers the question 'should we disable this service?'
func (svc *ServiceSpec) IsEnabled () (bool, error) {
	cmd := svc.BuildIsEnabledCommand()
	err := cmd.Run()
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if ok {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func (svc *ServiceSpec) IsExists () (bool, error) {
	// Build and run command. Non-ExitErrors are bad. ExitError.ErrorCode()=4 means service doesn't exist.
	cmd := svc.BuildDoesExistCommand()
	err := cmd.Run()
	if err != nil {
		exiterr, ok := err.(*exec.ExitError)
		if ok && exiterr.ExitCode() == 4 {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}

// Answers the question 'should we stop this service?'
func (svc *ServiceSpec) IsRunning () (bool, error) {
	cmd := svc.BuildIsActiveCommand()
	err := cmd.Run()
	if err != nil {
		_, ok := err.(*exec.ExitError)
		if ok {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}
func (svc *ServiceSpec) Remove (svcFS *ServiceFS) error {
	// Check if service exists. If so, try and disable
	is, err := svc.IsExists()
	if err != nil { return err }
	if is { 
		err := svc.Disable()
		if err != nil { return nil }
	}	

	// Remove service folder
	folderPath := svcFS.GetFolderPath()
	err = os.RemoveAll(folderPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			slog.Debug(fmt.Sprintf("%s already removed.", folderPath))
			return nil
		}
		return err
	}
	slog.Debug(fmt.Sprintf("%s successfully removed.", folderPath))

	return nil
}

func (svc *ServiceSpec) Start () error {
	// Check if service exists. If not, leave peacefully.
	is, err := svc.IsExists()
	slog.Debug(fmt.Sprintf("does=%t err=%s", is, err))
	if err != nil { return err }
	if !is { slog.Debug(fmt.Sprintf("%s doesn't exist.", svc.Name)); return nil }

	// Check if service is running. If so, leave peacefully. 
	is, err = svc.IsRunning()
	slog.Debug(fmt.Sprintf("is=%t err%s", is, err))
	if err != nil { return err; }
	if is { slog.Debug(fmt.Sprintf("%s is already running", svc.Name)); return nil }

	// Start service
	slog.Debug(fmt.Sprintf("Starting %s ...", svc.Name))
	cmd := svc.BuildStartCommand() 
	err = cmd.Run()
	if err != nil { return err }

	return nil
}

func (svc *ServiceSpec) Stop () error {
	// Check if service exists. If not, leave peacefully.
	is, err := svc.IsExists()
	slog.Debug(fmt.Sprintf("does=%t err=%s", is, err))
	if err != nil { return err }
	if !is { slog.Debug(fmt.Sprintf("%s doesn't exist.", svc.Name)); return nil }

	// Check if service is running. If so, stop it.
	is, err = svc.IsRunning()
	slog.Debug(fmt.Sprintf("is=%t err%s", is, err))
	if err != nil { return err; }
	if !is { slog.Debug(fmt.Sprintf("%s is already stopped.", svc.Name)); return nil }
	
	// Stop service
	slog.Debug(fmt.Sprintf("Stopping %s ...", svc.Name))
	cmd := svc.BuildStopCommand() 
	err = cmd.Run()
	if err != nil { return err }

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
type ServiceTemplateFS struct {
	fs.FS
}

func (fs *ServiceTemplateFS) GetRunScriptData(svcName string) (*TemplateData, error) {
	svcPattern := fmt.Sprintf("svc/%s/*", svcName)
	tmpl, err := template.ParseFS(fs.FS, svcPattern)
	if err != nil {
		return nil, err
	}
	runScriptFilename := "run.sh"
	tmpl = tmpl.Lookup(runScriptFilename)
	unitFileData := TemplateData{Template: *tmpl}

	return &unitFileData, nil
}

func (fs *ServiceTemplateFS) GetUnitFileData(svcName string) (*TemplateData, error) {
	svcPattern := fmt.Sprintf("svc/%s/*", svcName)
	tmpl, err := template.ParseFS(fs.FS, svcPattern)
	if err != nil {
		return nil, err
	}
	unitFilename := fmt.Sprintf("%s.service", svcName)
	tmpl = tmpl.Lookup(unitFilename)
	unitFileData := TemplateData{Template: *tmpl}

	return &unitFileData, nil
}
