package cmd

import (
	"strings"
	"testing"

	"github.com/dylt-dev/dylt/lib"
	"github.com/stretchr/testify/assert"
)

func TestRunCall0 (t *testing.T) {
	scriptPath := "/tmp/daylight.sh"
	scriptArgs := []string{"hello"}
	assert.FileExists(t, scriptPath)
	err := RunCall(scriptPath, scriptArgs)
	assert.Nil(t, err)
}

func TestRunCallCmd0 (t *testing.T) {
	sCmdline := "/Users/chris/src/dylt-dev/dylt/dylt call --script-path /tmp/daylight.sh hello"
	CheckRunCommandSuccess(sCmdline, t)
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
