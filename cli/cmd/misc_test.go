package cmd

import (
	"os"
	"strings"
	"testing"

	"text/template"

	"github.com/stretchr/testify/require"
)

func TestMisc (t *testing.T) {
	cmdName := "misc"
	cmdFlags := []string{}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, NewMiscCommand, cmdName, cmdFlags, cmdArgs, cmdString)
	require.IsType(t, &MiscCommand{}, cmd)
}



func TestCreateTwoNodeClusterCommand (t *testing.T) {
	// config get foo
	cmdName := "misc"
	subCmdName := "create-two-node-cluster"
	subCmdFlags := []string{}
	subCmdArgs := []string{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := NewMiscCommand(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	_TestSubcommandCreation[*CreateTwoNodeClusterCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
}

func TestGenEtcdRunScriptCommand  (t *testing.T) {
	// config get foo
	cmdName := "misc"
	subCmdName := "gen-etcd-run-script"
	subCmdFlags := []string{}
	subCmdArgs := []string{}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := NewMiscCommand(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	_TestSubcommandCreation[*GenEtcdRunScriptCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
}

func TestLookupCommand (t *testing.T) {
	// config get foo
	cmdName := "misc"
	subCmdName := "lookup"
	hostname := "hostname"
	subCmdFlags := []string{}
	subCmdArgs := []string{hostname}
	cmdline, cmdArgs, subCmdString := CreateCommandParams(cmdName, subCmdName, subCmdFlags, subCmdArgs)
	cmd := NewMiscCommand(cmdline, nil)
	// test parent command
	_TestParentCommand(t, cmd, cmdName, cmdArgs)
	// create + test subcommand
	subCmd := _TestSubcommandCreation[*LookupCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, hostname, subCmd.Hostname)
}

func TestGenEtcdRunScript(t *testing.T) {
	type EtcdRunScriptData struct {
		Name                string
		DataDir             string
		AdvertiseClientUrls []string
		ListenClientUrls    []string
		ClientCertAuth      bool
	}
	data := EtcdRunScriptData{
		Name:                "arleytown",
		AdvertiseClientUrls: []string{"https:/127.0.0.1:2239", "ip2:2239"},
		ListenClientUrls:    []string{"https:/127.0.0.1:2239"},
		ClientCertAuth:      false,
	}
	buf, err := content.ReadFile("content/run-etcd.sh.tmpl")
	require.NoError(t, err)
	tmpl := template.New("hello")
	tmpl.Funcs(template.FuncMap{
		"join": strings.Join,
	})
	tmpl, err = tmpl.Parse(string(buf))
	require.NoError(t, err)
	err = tmpl.Execute(os.Stdout, data)
	require.NoError(t, err)
}

func TestGetStdinStdoutStderrFdNums(t *testing.T) {
	var nStdin, nStdout, nStderr uintptr
	nStdin = os.Stdin.Fd()
	nStdout = os.Stdout.Fd()
	nStderr = os.Stderr.Fd()
	t.Logf("nStdin=%v", nStdin)
	t.Logf("nStdout=%v", nStdout)
	t.Logf("nStderr=%v", nStderr)
}

func TestNewCmdline (t *testing.T) {
	cmdName := "dylt"
	cmdFlags := []string{}
	cmdArgs := []string{}
	targetCmdline := Cmdline{"dylt"}
	cmdline := NewCmdline(cmdName, cmdFlags, cmdArgs)
	require.Equal(t, targetCmdline, cmdline)
}

func TestNewCmdlineArgs (t *testing.T) {
	cmdName := "dylt"
	cmdFlags := []string{}
	cmdArgs := []string{"config", "get", "foo"}
	targetCmdline := Cmdline{"dylt", "config", "get", "foo"}
	cmdline := NewCmdline(cmdName, cmdFlags, cmdArgs)
	require.Equal(t, targetCmdline, cmdline)
}

func TestNewCmdlineArgsFlags (t *testing.T) {
	cmdName := "dylt"
	cmdFlags := strings.Fields("--etcdDomain etcd.example.org")
	cmdArgs := []string{"get", "foo"}
	targetCmdline := Cmdline{"dylt", "--etcdDomain", "etcd.example.org", "get", "foo"}
	cmdline := NewCmdline(cmdName, cmdFlags, cmdArgs)
	require.Equal(t, targetCmdline, cmdline)
}

func TestNewCmdlineFlags (t *testing.T) {
	cmdName := "dylt"
	cmdFlags := strings.Fields("--etcdDomain etcd.example.org")
	cmdArgs := []string{}
	targetCmdline := Cmdline{"dylt", "--etcdDomain", "etcd.example.org"}
	cmdline := NewCmdline(cmdName, cmdFlags, cmdArgs)
	require.Equal(t, targetCmdline, cmdline)
}

func TestNewlineKiller(t *testing.T) {
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
	data := "MEAT!!!"
	PrintUsage(data)
}

func TestPrintMultiLineUsage_StringSlice(t *testing.T) {
	data := []string{"meat", "Meat", "MEAT"}
	PrintUsage(data)
}
