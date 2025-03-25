package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunInit (t *testing.T) {
	etcDomain := "hello.dylt.dev"
	err := RunInit(etcDomain)
	assert.Nil(t, err)
}

func TestInitCmd0 (t *testing.T) {
	if os.Getenv("DYLT_TEST_SYSTEST") == "" {
		t.Skip()
	}
	dyltPath, ok := os.LookupEnv("DYLT_EXE_PATH")
	assert.True(t, ok)
	assert.NotEmpty(t, dyltPath)
	sCmdline := fmt.Sprintf("%s init --etcd-domain hello.dylt.dev", dyltPath)
	err := CheckRunCommandSuccessNoOutput(sCmdline, t)
	assert.Nil(t, err)
}

func TestInitCmd1 (t *testing.T) {
	if os.Getenv("DYLT_TEST_SYSTEST") == "" {
		t.Skip()
	}
	dyltPath, ok := os.LookupEnv("DYLT_EXE_PATH")
	assert.True(t, ok)
	assert.NotEmpty(t, dyltPath)
	sCmdline := fmt.Sprintf("%s init --etcd-domain Hello-Hello-Hello.dylt.dev", dyltPath)
	err := CheckRunCommandSuccessNoOutput(sCmdline, t)
	assert.Nil(t, err)
}