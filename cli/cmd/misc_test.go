package cmd

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"text/template"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

func TestMisc (t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	cmdName := "misc"
	cmdFlags := []string{}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, MiscCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString)
	require.IsType(t, &BaseCommand[MiscOpts]{}, cmd)
}


func TestMiscHelp (t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	cmdName := "misc"
	cmdFlags := []string{"--help"}
	cmdArgs := []string{}
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t, MiscCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString).(*BaseCommand[MiscOpts])
	require.IsType(t, &BaseCommand[MiscOpts]{}, cmd)
	require.True(t, cmd.Help())
}


func TestMiscCreateTwoNodeClusterCommand (t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	// config get foo
	cmdName := "misc"
	subCmdName := "create-two-node-cluster"
	subCmdFlags := []string{}
	subCmdArgs := []string{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := MiscCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	_TestSubcommandCreation[*BaseCommand[CreateTwoNodeClusterOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
}


func TestMiscCreateTwoNodeClusterHelp(t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	cmdName := "misc"
	subCmdName := "create-two-node-cluster"
	subCmdFlags := []string{"--help"}
	subCmdArgs := []string{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := MiscCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test  subcommand
	subCmd := _TestSubcommandCreation[*BaseCommand[CreateTwoNodeClusterOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.True(t, subCmd.Help())
}


func TestMiscGenEtcdRunScript  (t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	// config get foo
	cmdName := "misc"
	subCmdName := "gen-etcd-run-script"
	subCmdFlags := []string{}
	subCmdArgs := []string{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := MiscCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	_TestSubcommandCreation[*BaseCommand[GenEtcdRunScriptOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
}

func TestMiscGenEtcdRunScriptHelp(t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	cmdName := "misc"
	subCmdName := "gen-etcd-run-script"
	subCmdFlags := []string{"--help"}
	subCmdArgs := []string{"foo"}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := MiscCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test  subcommand
	subCmd := _TestSubcommandCreation[*BaseCommand[GenEtcdRunScriptOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.True(t, subCmd.Help())
}

func TestMiscLookup (t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	// config get foo
	cmdName := "misc"
	subCmdName := "lookup"
	hostname := "hostname"
	subCmdFlags := []string{}
	subCmdArgs := []string{hostname}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := MiscCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	subCmd := _TestSubcommandCreation[*BaseCommand[LookupOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, hostname, subCmd.opts.Hostname)
}


func TestMiscLookupHelp(t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	cmdName := "misc"
	subCmdName := "lookup"
	subCmdFlags := []string{"--help"}
	subCmdArgs := []string{"hostname"}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := MiscCommandF.New(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test  subcommand
	subCmd := _TestSubcommandCreation[*BaseCommand[LookupOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.True(t, subCmd.Help())
}

func TestGetStdinStdoutStderrFdNums(t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	var nStdin, nStdout, nStderr uintptr
	nStdin = os.Stdin.Fd()
	nStdout = os.Stdout.Fd()
	nStderr = os.Stderr.Fd()
	t.Logf("nStdin=%v", nStdin)
	t.Logf("nStdout=%v", nStdout)
	t.Logf("nStderr=%v", nStderr)
}

func TestNewCmdline (t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	cmdName := "dylt"
	cmdFlags := []string{}
	cmdArgs := []string{}
	targetCmdline := Cmdline{"dylt"}
	cmdline := NewCmdline(cmdName, cmdFlags, cmdArgs)
	require.Equal(t, targetCmdline, cmdline)
}

func TestNewCmdlineArgs (t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	cmdName := "dylt"
	cmdFlags := []string{}
	cmdArgs := []string{"config", "get", "foo"}
	targetCmdline := Cmdline{"dylt", "config", "get", "foo"}
	cmdline := NewCmdline(cmdName, cmdFlags, cmdArgs)
	require.Equal(t, targetCmdline, cmdline)
}

func TestNewCmdlineArgsFlags (t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	cmdName := "dylt"
	cmdFlags := strings.Fields("--etcdDomain etcd.example.org")
	cmdArgs := []string{"get", "foo"}
	targetCmdline := Cmdline{"dylt", "--etcdDomain", "etcd.example.org", "get", "foo"}
	cmdline := NewCmdline(cmdName, cmdFlags, cmdArgs)
	require.Equal(t, targetCmdline, cmdline)
}

func TestNewCmdlineFlags (t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	cmdName := "dylt"
	cmdFlags := strings.Fields("--etcdDomain etcd.example.org")
	cmdArgs := []string{}
	targetCmdline := Cmdline{"dylt", "--etcdDomain", "etcd.example.org"}
	cmdline := NewCmdline(cmdName, cmdFlags, cmdArgs)
	require.Equal(t, targetCmdline, cmdline)
}

func TestNewlineKiller(t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	s := `
	line1
	{{- if .Line2 }}
	{{ .Line2 -}}
	{{ end }}
	line3
	`

	type nkdata struct{ Line2 string }
	tmpl := template.New("test")
	tmpl, err := tmpl.Parse(s)
	require.NoError(t, err)
	err = tmpl.Execute(os.Stdout, nkdata{})
	require.NoError(t, err)
	err = tmpl.Execute(os.Stdout, nkdata{Line2: "MEAT"})
	require.NoError(t, err)
}

func TestPrintMultiLineUsage_String(t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	data := "MEAT!!!"
	PrintUsage(data)
}

func TestPrintMultiLineUsage_StringSlice(t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	data := []string{"meat", "Meat", "MEAT"}
	PrintUsage(data)
}
