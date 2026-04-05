package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHelpFlag (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cfg := BaseCommandConfig[EmptyOpts]{
		name: "base",
		usage: "",
		validator: ArgCountValidator{nExpected: 2},
	}
	cmd := NewBaseCommand([]string{"dylt", "--help"}, nil, cfg)
	err := cmd.HandleArgs()
	fmt.Printf("cmd.Help()=%v\n", cmd.Help())
	require.NoError(t, err)
	require.True(t, cmd.Help())
}