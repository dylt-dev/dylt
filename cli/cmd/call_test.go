package cmd

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/dylt-dev/dylt/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCall (t *testing.T) {
	cmdName := "call"
	cmdFlags := []string{}
	cmdArgs := []string{"command", "foo", "bar", "bum"}
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t, NewCallCommand, cmdName, cmdFlags, cmdArgs, cmdString)
	require.IsType(t, &CallCommand{}, cmd)
}


func TestCallHelp(t *testing.T) {
	cmdName := "call"
	cmdFlags := []string{"--help"}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, NewCallCommand, cmdName, cmdFlags, cmdArgs, cmdString)
	require.True(t, cmd.Help)
}

func TestRunCall0 (t *testing.T) {
	scriptPath := "/tmp/daylight.sh"
	_, err := os.Stat(scriptPath)
	if os.IsNotExist(err) {
		t.Skipf("script not found: %s", scriptPath)
	}
	scriptArgs := []string{"hello"}
	assert.FileExists(t, scriptPath)
	err = RunCall(scriptPath, scriptArgs)
	assert.Nil(t, err)
}

func TestRunCallCmd0 (t *testing.T) {
	sCmdline := "/Users/chris/src/dylt-dev/dylt/dylt call --script-path /tmp/daylight.sh hello"
	lib.CheckRunCommandSuccess(sCmdline, t)
}

func TestCallNoScript (t *testing.T) {
	sCmdline := "XXX call --script-path /tmp/daylight.sh hello"
	var cmdline Cmdline = strings.Split(sCmdline, " ")
	t.Log("cmdline.Command()", cmdline.Command(), "cmdline.Args()", cmdline.Args())
	rc, stdout, err := lib.RunCommand(cmdline.Command(), cmdline.Args()...)
	assert.NotEqual(t, 0, rc)
	assert.Empty(t, stdout)
	assert.NotNil(t, err)
	t.Log(string(stdout))
}
