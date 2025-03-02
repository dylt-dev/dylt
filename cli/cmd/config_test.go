package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunConfigGet (t *testing.T) {
	key := "etcd-domain"
	err := RunConfigGet(key)
	assert.Nil(t, err)
}

func TestConfigGetCmd (t *testing.T) {
	sCmdline := "/Users/chris/src/dylt-dev/dylt/dylt config get etcd-domain"
	CheckRunCommandSuccess(sCmdline, t)
}


func TestRunConfigSet (t *testing.T) {
	key := "etcd-domain"
	val := "poo"
	err := RunConfigSet(key, val)
	assert.Nil(t, err)
}

func TestConfigSetCmd0 (t *testing.T) {
	sCmdline := "/Users/chris/src/dylt-dev/dylt/dylt config set etcd-domain MOO"
	CheckRunCommandSuccessNoOutput(sCmdline, t)
}

func TestConfigSetCmd1 (t *testing.T) {
	sCmdline := "/Users/chris/src/dylt-dev/dylt/dylt config set etcd-domain hello.dylt.dev"
	CheckRunCommandSuccessNoOutput(sCmdline, t)
}


func TestRunConfigShow (t *testing.T) {
	err := RunConfigShow()
	assert.Nil(t, err)
}

func TestConfigShowCmd (t *testing.T) {
	sCmdline := "/Users/chris/src/dylt-dev/dylt/dylt config show"
	CheckRunCommandSuccess(sCmdline, t)
}
