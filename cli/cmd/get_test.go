package cmd

import (
	"fmt"
	"testing"

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

