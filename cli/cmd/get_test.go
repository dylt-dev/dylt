package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet (t *testing.T) {
	cmdName := "get"
	key := "foo"
	cmdFlags := []string{}
	cmdArgs := []string{key}
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t, NewGetCommand, cmdName, cmdFlags, cmdArgs, cmdString)
	require.IsType(t, &GetCommand{}, cmd)
	require.Equal(t, key, cmd.Key)
}


func TestRunGet (t *testing.T) {
	key := "hello"
	err := RunGet(key)
	assert.Nil(t, err)
}

func TestGetCmd0 (t *testing.T) {
	if os.Getenv("DYLT_TEST_SYSTEST") == "" {
		t.Skip()
	}
	dyltPath, ok := os.LookupEnv("DYLT_EXE_PATH")
	assert.True(t, ok)
	assert.NotEmpty(t, dyltPath)
	sCmdline := fmt.Sprintf("%s get hello", dyltPath)
	err := CheckRunCommandSuccess(sCmdline, t)
	assert.Nil(t, err)
}
