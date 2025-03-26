package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/dylt-dev/dylt/lib"
	"github.com/dylt-dev/dylt/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	fsSvcFiles, err := fs.Sub(lib.EMBED_SvcFiles, "svcfiles")
	assert.NoError(t, err)
	tmplSvc, err := template.GetServiceTemplate(fsSvcFiles, "watch-daylight")
	assert.NoError(t, err)
	assert.NotNil(t, tmplSvc)
	tmpl := tmplSvc.Lookup("/run.sh")
	require.NotNil(t, tmpl)
	err = tmpl.Execute(os.Stdout, map[any]any{})
	assert.NoError(t, err)
}

func TestEmitWatchDaylightUnitFile(t *testing.T) {
	fsSvcFiles, err := fs.Sub(lib.EMBED_SvcFiles, "svcfiles")
	assert.NoError(t, err)
	tmplSvc, err := template.GetServiceTemplate(fsSvcFiles, "watch-daylight")
	assert.NoError(t, err)
	assert.NotNil(t, tmplSvc)
	tmpl := tmplSvc.Lookup("/watch-daylight.service")
	require.NotNil(t, tmpl)
	err = tmpl.Execute(os.Stdout, map[any]any{})
	assert.NoError(t, err)
}

func TestChmodR0(t *testing.T) {
	svcPath := "/opt/svc/watch-daylight"
	// uid + gid for local user on local workstation
	err := lib.ChownR(svcPath, 501, 20)
	assert.Nil(t, err)
}

func Test_WatchDaylight_WriteRunScript(t *testing.T) {
	fsSvcFiles, err := fs.Sub(lib.EMBED_SvcFiles, "svcfiles")
	assert.NoError(t, err)
	svcName := "watch-daylight"
	tmpl, err := template.GetServiceTemplate(fsSvcFiles, svcName)
	assert.NoError(t, err)
	assert.NotNil(t, tmpl)
	svc := lib.NewServiceSpec(svcName)
	err = tmpl.WriteRunScript(os.Stdout, svc.Data)
	assert.NoError(t, err)
}

func Test_WatchDaylight_WriteUnitFile(t *testing.T) {
	fsSvcFiles, err := fs.Sub(lib.EMBED_SvcFiles, "svcfiles")
	assert.NoError(t, err)
	svcName := "watch-daylight"
	tmpl, err := template.GetServiceTemplate(fsSvcFiles, svcName)
	assert.NoError(t, err)
	assert.NotNil(t, tmpl)
	svc := lib.NewServiceSpec(svcName)
	err = tmpl.WriteUnitFile(os.Stdout, svc.Data)
	assert.NoError(t, err)
}
