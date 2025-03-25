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

func TestRunHost(t *testing.T) {
	if os.Getenv("DYLT_TEST_SYSTEST") == "" {
		t.Skip()
	}
	cmdName := "init"
	cmdArgs := []string{}
	err := RunHost(cmdName, cmdArgs)
	assert.Nil(t, err)
}

func TestHostCmd0(t *testing.T) {
	if os.Getenv("DYLT_TEST_SYSTEST") == "" {
		t.Skip()
	}
	dyltPath, ok := os.LookupEnv("DYLT_EXE_PATH")
	assert.True(t, ok)
	assert.NotEmpty(t, dyltPath)
	cmd := fmt.Sprintf("%s host", dyltPath)
	err := CheckRunCommandSuccess(cmd, t)
	assert.Nil(t, err)
}

func TestRunHostInit(t *testing.T) {
	if os.Getenv("DYLT_TEST_SYSTEST") == "" {
		t.Skip()
	}
	err := RunHostInit(501, 20)
	assert.Nil(t, err)
}

func TestHostInitCmd0(t *testing.T) {
	if os.Getenv("DYLT_TEST_SYSTEST") == "" {
		t.Skip()
	}
	dyltPath, ok := os.LookupEnv("DYLT_EXE_PATH")
	assert.True(t, ok)
	assert.NotEmpty(t, dyltPath)
	cmd := fmt.Sprintf("%s host init", dyltPath)
	err := CheckRunCommandSuccess(cmd, t)
	assert.Nil(t, err)
}

// not-a-test
// Print out all the files in EMBED_SvcFiles. Useful sanity check.
func TestWalkSvcFolder(t *testing.T) {
	fs.WalkDir(lib.EMBED_SvcFiles, ".", func(p string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			fmt.Printf("%s\n", p)
		}
		return nil
	})
}

func TestEmitWatchDaylightRunScript(t *testing.T) {
	// assert.True(t, lib.PATH_WatchDaylightRunScript)
	tmpl, err := template.ParseFS(lib.EMBED_SvcFiles, "svc/watch-daylight/*")
	assert.Nil(t, err)
	assert.NotNil(t, tmpl)
	tmpl = tmpl.Lookup("run.sh")
	assert.NotNil(t, tmpl)
	err = tmpl.Execute(os.Stdout, map[any]any{})
	assert.Nil(t, err)
}

func TestEmitWatchDaylightUnitFile(t *testing.T) {
	// assert.True(t, lib.PATH_WatchDaylightRunScript)
	tmpl, err := template.ParseFS(lib.EMBED_SvcFiles, "svc/watch-daylight/*")
	assert.Nil(t, err)
	assert.NotNil(t, tmpl)
	tmpl = tmpl.Lookup("watch-daylight.service")
	assert.NotNil(t, tmpl)
	data := map[string]string{}
	err = tmpl.Execute(os.Stdout, data)
	assert.Nil(t, err)
}

func TestChmodR0(t *testing.T) {
	svcPath := "/opt/svc/watch-daylight"
	// uid + gid for local user on local workstation
	err := lib.ChownR(svcPath, 501, 20)
	assert.Nil(t, err)
}

func Test_WatchDaylight_WriteRunScript(t *testing.T) {
	svcName := "watch-daylight"
	svc := lib.NewServiceSpec(svcName)
	svcFS := lib.NewServiceFS(svcName, "/opt/svc")
	templateFS := lib.ServiceTemplateFS{FS: lib.EMBED_SvcFiles}
	err := svcFS.WriteRunScript(svc, &templateFS)
	assert.Nil(t, err)
}

func Test_WatchDaylight_WriteUnitFile(t *testing.T) {
	svcName := "watch-daylight"
	svc := lib.NewServiceSpec(svcName)
	svcFS := lib.NewServiceFS(svcName, "/opt/svc")
	templateFS := lib.ServiceTemplateFS{FS: lib.EMBED_SvcFiles}
	err := svcFS.WriteUnitFile(svc, &templateFS)
	assert.Nil(t, err)
}
