package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHelpFlag (t *testing.T) {
	cmd := NewBaseCommand("base", []string{"dylt", "--help"}, nil, "usage", nil, ArgCountValidator{nExpected: 3})
	err := cmd.HandleArgs()
	fmt.Printf("cmd.Help=%v\n", cmd.Help)
	require.NoError(t, err)
	require.True(t, cmd.Help) 
}