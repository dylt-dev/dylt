package cmd

import (
	"fmt"
	"testing"

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
	err := RunHostInit()
	assert.Nil(t, err)
}

func TestHostInitCmd0 (t *testing.T) {
	cmd := fmt.Sprintf("%s host init", PATH_Dylt)
	err := CheckRunCommandSuccess(cmd, t)
	assert.Nil(t, err)
}
