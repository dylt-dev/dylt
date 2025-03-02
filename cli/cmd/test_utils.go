package cmd

import (
	"strings"
	"testing"

	"github.com/dylt-dev/dylt/lib"
	"github.com/stretchr/testify/assert"
)

const PATH_Dylt = "/Users/chris/src/dylt-dev/dylt/dylt"

func CheckRunCommandSuccess (sCmdline string, t *testing.T) error {
	var cmdline Cmdline = strings.Split(sCmdline, " ")
	rc, stdout, err := lib.RunCommand(cmdline.Command(), cmdline.Args()...)
	assert.Equal(t, 0, rc)
	assert.NotEmpty(t, stdout)
	assert.Nil(t, err)
	t.Log(string(stdout))
	return err
}


func CheckRunCommandSuccessNoOutput (sCmdline string, t *testing.T) error {
	var cmdline Cmdline = strings.Split(sCmdline, " ")
	rc, stdout, err := lib.RunCommand(cmdline.Command(), cmdline.Args()...)
	assert.Equal(t, 0, rc)
	assert.Empty(t, stdout)
	assert.Nil(t, err)
	t.Log(string(stdout))
	return err
}