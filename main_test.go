package main

import (
	"strings"
	"testing"

	clicmd "github.com/dylt-dev/dylt/cli/cmd"
	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

// cmdline: dylt vm list
// Expected: --help=False SubCommand="vm" SubArgs="list")
func TestVmMainList(t *testing.T) {
	cmdline := []string{"vm", "list"}
	cmd := clicmd.MainCommandF.New(cmdline, nil).(*clicmd.BaseCommand[clicmd.MainOpts])
	cmd.HandleArgs()
	require.False(t, cmd.Help())
	require.Equal(t, "vm", cmd.SubCommand)
	require.Equal(t, []string{"list"}, cmd.SubArgs)
}

func TestRun(t *testing.T) {
	testCommandLine(t, "")
}

// dylt --help
func TestRun_Help(t *testing.T) {
	testCommandLine(t, "--help")
}

func TestRun_Status(t *testing.T) {
	testCommandLine(t, "status")
}

func TestRun_StatusHelp(t *testing.T) {
	common.InitLogging()
	var cmdline clicmd.Cmdline = []string{"status", "--help"}
	cmd := clicmd.MainCommandF.New(cmdline, nil)
	cmd.Run()
}

func testCommandLine(t *testing.T, s string) {
	common.InitLogging()
	var cmdline clicmd.Cmdline = strings.Fields(s)
	cmd := clicmd.MainCommandF.New(cmdline, nil)
	err := cmd.Run()
	require.NoError(t, err)
}