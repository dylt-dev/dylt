package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunInit (t *testing.T) {
	etcDomain := "hello.dylt.dev"
	err := RunInit(etcDomain)
	assert.Nil(t, err)
}

func TestInitCmd0 (t *testing.T) {
	sCmdline := fmt.Sprintf("%s init --etcd-domain hello.dylt.dev", PATH_Dylt)
	err := CheckRunCommandSuccessNoOutput(sCmdline, t)
	assert.Nil(t, err)
}

func TestInitCmd1 (t *testing.T) {
	sCmdline := fmt.Sprintf("%s init --etcd-domain Hello-Hello-Hello.dylt.dev", PATH_Dylt)
	err := CheckRunCommandSuccessNoOutput(sCmdline, t)
	assert.Nil(t, err)
}