package cmd

import (
	"fmt"
	"testing"

	"github.com/dylt-dev/dylt/common"
	"github.com/dylt-dev/dylt/api"
	"github.com/stretchr/testify/require"
)

func TestHelpFlag (t *testing.T) {
	fnTeardown := common.Setup(t)
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

// Legacy tagless fallback examples — reference patterns for constructing
// commands without struct tags. These mirror what the tag system automates.

func TestUntaggedListCommand(t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)

	cmdline := NewCmdline("list", nil, nil)
	cfg := BaseCommandConfig[ListOpts]{
		name:      "list",
		opts:      ListOpts{},
		fnRun:     func(cmd *BaseCommand[ListOpts]) error { return api.RunList() },
		usage:     CreateUsageString(USG_List),
		validator: ArgCountValidator{nExpected: 0},
	}
	cmd := NewBaseCommand(cmdline, nil, cfg)

	err := cmd.HandleArgs()
	require.NoError(t, err)
	require.Equal(t, "list", cmd.CommandName())
}

func TestUntaggedVmSetCommand(t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)

	cmdline := NewCmdline("vm.set", nil, []string{"myvm", "cpu", "4"})
	fnRun := func(cmd *BaseCommand[VmSetOpts]) error {
		return api.RunVmSet(cmd.opts.Name, cmd.opts.Key, cmd.opts.Value)
	}
	cfg := BaseCommandConfig[VmSetOpts]{
		name:      "vm.set",
		fnRun:     fnRun,
		opts:      VmSetOpts{},
		usage:     CreateUsageString(USG_Vm_Set),
		validator: ArgCountValidator{nExpected: 3},
	}
	cmd := NewBaseCommand(cmdline, nil, cfg)
	cmd.argMap = ArgMap{
		0: func(s string) { cmd.opts.Name = s },
		1: func(s string) { cmd.opts.Key = s },
		2: func(s string) { cmd.opts.Value = s },
	}

	err := cmd.HandleArgs()
	require.NoError(t, err)
	require.Equal(t, "myvm", cmd.opts.Name)
	require.Equal(t, "cpu", cmd.opts.Key)
	require.Equal(t, "4", cmd.opts.Value)
}

type testFlagPosOpts struct {
	Name    string
	Timeout string
	Mode    string
	Verbose bool
}

func TestUntaggedCompositeCommand(t *testing.T) {
	fnTeardown := common.Setup(t)
	defer fnTeardown(t)

	cmdline := NewCmdline("example", []string{"--verbose", "--mode", "strict"}, []string{"myvm", "30s"})
	fnRun := func(cmd *BaseCommand[testFlagPosOpts]) error { return nil }
	cfg := BaseCommandConfig[testFlagPosOpts]{
		name:      "example",
		opts:      testFlagPosOpts{},
		fnRun:     fnRun,
		usage:     "usage info",
		validator: ArgCountValidator{nExpected: 2},
	}
	cmd := NewBaseCommand(cmdline, nil, cfg)
	cmd.argMap = ArgMap{
		0: func(s string) { cmd.opts.Name = s },
		1: func(s string) { cmd.opts.Timeout = s },
	}
	cmd.StringVar(&cmd.opts.Mode, "mode", "default", "operating mode")
	cmd.BoolVar(&cmd.opts.Verbose, "verbose", false, "enable verbose output")

	err := cmd.HandleArgs()
	require.NoError(t, err)
	require.Equal(t, "myvm", cmd.opts.Name)
	require.Equal(t, "30s", cmd.opts.Timeout)
	require.Equal(t, "strict", cmd.opts.Mode)
	require.True(t, cmd.opts.Verbose)
}