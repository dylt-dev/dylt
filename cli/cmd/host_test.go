package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/dylt-dev/dylt/api"
	"github.com/dylt-dev/dylt/common"
	"github.com/dylt-dev/dylt/lib"
	"github.com/dylt-dev/dylt/systemd"
	"github.com/dylt-dev/dylt/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHost(t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)

	cmdName := "host"
	cmdFlags := []string{}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, HostCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString).(*BaseCommand[HostOpts])
	require.NotNil(t, cmd)
}

func TestHostHelp(t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)

	cmdName := "host"
	cmdFlags := []string{"--help"}
	cmdArgs := []string{}
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t, HostCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString).(*BaseCommand[HostOpts])
	require.NotNil(t, cmd)
	require.True(t, cmd.Help())
}

func TestHostInit(t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)

	// config get foo
	cmdName := "host"
	subCmdName := "init"
	gid := 1000
	uid := 2000
	subCmdFlags := []string{"--gid", fmt.Sprint(gid), "--uid", fmt.Sprint(uid)}
	subCmdArgs := []string{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := HostCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	subCmd := _TestSubcommandCreation[*BaseCommand[HostInitOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, gid, subCmd.opts.Gid)
	require.Equal(t, uid, subCmd.opts.Uid)
}

func TestHostInitHelp(t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)

	cmdName := "host"
	subCmdName := "init"
	subCmdFlags := []string{"--help"}
	subCmdArgs := []string{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := HostCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test  subcommand
	subCmd := _TestSubcommandCreation[*BaseCommand[HostInitOpts]](
		t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.True(t, subCmd.Help())
}

func TestRunHost(t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)

	if os.Getenv("DYLT_TEST_SYSTEST") == "" {
		t.Skip()
	}
	cmdName := "init"
	cmdArgs := []string{}
	cmdline := append(Cmdline{cmdName}, cmdArgs...)
	parent := Command(nil)
	err := RunHost(cmdline, parent)
	assert.Nil(t, err)
}

func TestHostCmd0(t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)

	if os.Getenv("DYLT_TEST_SYSTEST") == "" {
		t.Skip()
	}
	dyltPath, ok := os.LookupEnv("DYLT_EXE_PATH")
	assert.True(t, ok)
	assert.NotEmpty(t, dyltPath)
	cmd := fmt.Sprintf("%s host", dyltPath)
	err := lib.CheckRunCommandSuccess(cmd, t)
	assert.Nil(t, err)
}

func TestEmitWatchDaylightRunScript(t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)

	fsSvcFiles, err := fs.Sub(api.EMBED_SvcFiles, "svcfiles")
	assert.NoError(t, err)
	tmplSvc, err := template.NewTemplate(fsSvcFiles, "watch-daylight")
	assert.NoError(t, err)
	assert.NotNil(t, tmplSvc)
	tmpl := tmplSvc.Lookup("/run.sh")
	require.NotNil(t, tmpl)
	err = tmpl.Execute(os.Stdout, map[any]any{})
	assert.NoError(t, err)
}

func TestEmitWatchDaylightUnitFile(t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)

	fsSvcFiles, err := fs.Sub(api.EMBED_SvcFiles, "svcfiles")
	assert.NoError(t, err)
	tmplSvc, err := template.NewTemplate(fsSvcFiles, "watch-daylight")
	assert.NoError(t, err)
	assert.NotNil(t, tmplSvc)
	tmpl := tmplSvc.Lookup("/watch-daylight.service")
	require.NotNil(t, tmpl)
	err = tmpl.Execute(os.Stdout, map[any]any{})
	assert.NoError(t, err)
}

func TestChmodR0(t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)

	svcPath := "/opt/svc/watch-daylight"
	// uid + gid for local user on local workstation
	err := common.ChownR(svcPath, 501, 20)
	assert.Nil(t, err)
}

func Test_WatchDaylight_WriteRunScript(t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)

	fsSvcFiles, err := fs.Sub(api.EMBED_SvcFiles, "svcfiles")
	assert.NoError(t, err)
	svcName := "watch-daylight"
	tmpl, err := template.NewTemplate(fsSvcFiles, svcName)
	assert.NoError(t, err)
	assert.NotNil(t, tmpl)
	svc := systemd.NewServiceSpec(svcName)
	err = tmpl.WriteRunScript(os.Stdout, svc.Data)
	assert.NoError(t, err)
}

func Test_WatchDaylight_WriteUnitFile(t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)

	fsSvcFiles, err := fs.Sub(api.EMBED_SvcFiles, "svcfiles")
	assert.NoError(t, err)
	svcName := "watch-daylight"
	tmpl, err := template.NewTemplate(fsSvcFiles, svcName)
	assert.NoError(t, err)
	assert.NotNil(t, tmpl)
	svc := systemd.NewServiceSpec(svcName)
	err = tmpl.WriteUnitFile(os.Stdout, svc.Data)
	assert.NoError(t, err)
}
