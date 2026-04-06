package cmd

import (
	"fmt"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/stretchr/testify/require"
)

func TestGet (t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	cmdName := "get"
	key := "foo"
	cmdFlags := []string{}
	cmdArgs := []string{key}
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t, GetCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString)
	require.Equal(t, key, *cmd.ArgMap()[0])
}


func TestGetHelp (t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)
	
	cmdName := "get"
	cmdFlags := []string{"--help"}
	cmdArgs := []string{}
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t, GetCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString)
	require.True(t, cmd.Help())
}