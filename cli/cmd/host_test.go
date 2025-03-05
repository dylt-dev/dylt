package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"testing"
	"text/template"

	"github.com/dylt-dev/dylt/lib"
	"github.com/stretchr/testify/assert"
)

func TestRunHost (t *testing.T) {
	cmdName := "init"
	cmdArgs := []string{}
	err := RunHost(cmdName, cmdArgs)
	assert.Nil(t, err)
}

func TestHostCmd0 (t *testing.T) {
	cmd := fmt.Sprintf("%s host", PATH_Dylt)
	err := CheckRunCommandSuccess(cmd, t)
	assert.Nil(t, err)
}

func TestRunHostInit (t *testing.T) {
	err := RunHostInit(501, 20)
	assert.Nil(t, err)
}

func TestHostInitCmd0 (t *testing.T) {
	cmd := fmt.Sprintf("%s host init", PATH_Dylt)
	err := CheckRunCommandSuccess(cmd, t)
	assert.Nil(t, err)
}

func TestWalkSvcFolder (t *testing.T) {
	fs.WalkDir(lib.FOL_Svc, ".", func(p string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			fmt.Printf("%s\n", p)
		}
		return nil
	})
}

func TestEmitWatchDaylightRunScript(t *testing.T) {
	// assert.True(t, lib.PATH_WatchDaylightRunScript)
	tmpl, err := template.ParseFS(lib.FOL_Svc, "svc/watch-daylight/*")
	assert.Nil(t, err)
	assert.NotNil(t, tmpl)
	tmpl = tmpl.Lookup(lib.FN_WatchDaylightRunScript)
	assert.NotNil(t, tmpl)
	data := map[string]string{}
	err = tmpl.Execute(os.Stdout, data)
	assert.Nil(t, err)
}

func TestEmitWatchDaylightUnitFile (t *testing.T) {
	// assert.True(t, lib.PATH_WatchDaylightRunScript)
	tmpl, err := template.ParseFS(lib.FOL_Svc, "svc/watch-daylight/*")
	assert.Nil(t, err)
	assert.NotNil(t, tmpl)
	tmpl = tmpl.Lookup(lib.FN_WatchDaylightUnitFile)
	assert.NotNil(t, tmpl)
	data := map[string]string{}
	err = tmpl.Execute(os.Stdout, data)
	assert.Nil(t, err)
}

func TestChmodR0 (t *testing.T) {
	svcPath := "/opt/svc/watch-daylight-go"
	// uid + gid for local user on local workstation
	err := lib.ChownR(svcPath, 501, 20)
	assert.Nil(t, err)
}

func Test_WatchDaylight_WriteRunScript (t *testing.T) {
	svcName := "watch-daylight"
	err := lib.WriteRunScript(svcName, lib.ServiceData{})
	assert.Nil(t, err)
}

func Test_WatchDaylight_WriteUnitFile (t *testing.T) {
	svcName := "watch-daylight"
	err := lib.WriteUnitFile(svcName, lib.ServiceData{})
	assert.Nil(t, err)
}
