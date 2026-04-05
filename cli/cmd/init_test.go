package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/dylt-dev/dylt/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdName := "init"
	etcdDomain := "foo.dylt.dev"
	cmdFlags := []string{"--etcd-domain", etcdDomain}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, InitCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString).(*BaseCommand[InitOpts])
	require.NotNil(t, cmd)
	require.Equal(t, etcdDomain, cmd.opts.EtcdDomain)
	
}

func TestInitHelp(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdName := "init"
	cmdFlags := []string{"--help"}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, InitCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString).(*BaseCommand[InitOpts])
	require.NotNil(t, cmd)
	require.True(t, cmd.Help())
}

func TestRunInit (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	etcDomain := "hello.dylt.dev"
	err := api.RunInit(etcDomain)
	assert.Nil(t, err)
}

func TestInitCmd0 (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
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
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
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