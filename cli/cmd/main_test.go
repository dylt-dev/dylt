package cmd

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Returns
//     cmdline	    The command line for the command+subcommand invocation
//     cmdArgs      The arguments for the command
//     subCmdString The command string for the full command: cmd, subCmd, subCmdArgs
func CreateCommandParams(cmdName string,
	                     subCmdName string,
						 subCmdFlags []string,
						 subCmdArgs []string) (cmdline Cmdline, cmdArgs Cmdline, cmdString string) {
	subCmdline := NewCmdline(subCmdName, subCmdFlags, subCmdArgs)
	cmdline = NewCmdline(cmdName, []string{}, subCmdline)
	cmdArgs = cmdline.Args()
	cmdString = strings.Join(cmdline[0:2], " ")
	return
}

func TestMain(t *testing.T) {
	cmdName := "dylt"
	cmdFlags := []string{}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, NewMainCommand, cmdName, cmdFlags, cmdArgs, cmdString)
	require.False(t, cmd.Help)
}

func TestMainHelp(t *testing.T) {
	cmdName := "dylt"
	cmdFlags := []string{"--help"}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd := CreateAndTestCommand(t, NewMainCommand, cmdName, cmdFlags, cmdArgs, cmdString)
	require.True(t, cmd.Help)
}

// dylt call foo
func TestMainSubCall(t *testing.T) {
	// Independent values
	cmdName := "dylt"
	cmdFlags := []string{}
	subCmdName := "call"
	subCmdFlags := []string{}
	subCmdArgs := []string{"foo"}

	// Create dependent values for command + test
	cmdArgs := slices.Concat([]string{subCmdName}, subCmdFlags, subCmdArgs)
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t,
		NewMainCommand,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

		// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*CallCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &CallCommand{}, subCmd)
}

func TestMainSubConfig (t *testing.T) {
	// Independent values
	cmdName := "dylt"
	cmdFlags := []string{}
	subCmdName := "config"
	subCmdFlags := []string{}
	subCmdArgs := []string{}

	// Create dependent values for command + test
	cmdArgs := slices.Concat([]string{subCmdName}, subCmdFlags, subCmdArgs)
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t,
		NewMainCommand,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

		// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*ConfigCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &ConfigCommand{}, subCmd)
}


// dylt config get foo
func TestMainSubConfigGetFoo(t *testing.T) {
	// Independent values
	cmdName := "dylt"
	cmdFlags := []string{}
	subCmdName := "config"
	subCmdFlags := []string{}
	subCmdArgs := []string{"get", "foo"}

	// Create dependent values for command + test
	cmdArgs := slices.Concat([]string{subCmdName}, subCmdFlags, subCmdArgs)
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t,
		NewMainCommand,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

		// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*ConfigCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &ConfigCommand{}, subCmd)
}

// dylt get foo
func TestMainSubGet(t *testing.T) {
	// Independent values
	cmdName := "dylt"
	cmdFlags := []string{}
	subCmdName := "get"
	subCmdFlags := []string{}
	subCmdArgs := []string{"foo"}

	// Create dependent values for command + test
	cmdArgs := slices.Concat([]string{subCmdName}, subCmdFlags, subCmdArgs)
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t,
		NewMainCommand,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

		// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*GetCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &GetCommand{}, subCmd)
}


func TestMainSubHost(t *testing.T) {
	// Independent values
	cmdName := "dylt"
	cmdFlags := []string{}
	subCmdName := "host"
	subCmdFlags := []string{}
	subCmdArgs := []string{}

	// Create dependent values for command + test
	cmdArgs := slices.Concat([]string{subCmdName}, subCmdFlags, subCmdArgs)
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t,
		NewMainCommand,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

	// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*HostCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &HostCommand{}, subCmd)
}


func TestMainSubHostInit(t *testing.T) {
	// Independent values
	cmdName := "dylt"
	cmdFlags := []string{}
	subCmdName := "host"
	gid := 1000
	uid := 1000
	subCmdFlags := []string{}
	subCmdArgs := []string{"init", "--gid", fmt.Sprint(gid), "--uid", fmt.Sprint(uid)}

	// Create dependent values for command + test
	cmdArgs := slices.Concat([]string{subCmdName}, subCmdFlags, subCmdArgs)
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t,
		NewMainCommand,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

	// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*HostCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &HostCommand{}, subCmd)
}


// dylt init --etcd-domain foo.dylt.dev
func TestMainSubInit(t *testing.T) {
	// Create + parse main command with "dylt get foo"
	// Independent values
	cmdName := "dylt"
	cmdFlags := []string{}
	subCmdName := "init"
	etcdDomain := "foo.dylt.dev"
	subCmdFlags := []string{"--etcd-domain", etcdDomain}
	subCmdArgs := []string{}

	// Create dependent values for command + test
	cmdArgs := slices.Concat([]string{subCmdName}, subCmdFlags, subCmdArgs)
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t,
		NewMainCommand,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

	// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*InitCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, etcdDomain, subCmd.EtcdDomain)
}


// dylt list 
func TestMainSubList(t *testing.T) {
	// Independent values
	cmdName := "dylt"
	cmdFlags := []string{}
	subCmdName := "list"
	subCmdFlags := []string{}
	subCmdArgs := []string{}

	// Create dependent values for command + test
	cmdArgs := slices.Concat([]string{subCmdName}, subCmdFlags, subCmdArgs)
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t,
		NewMainCommand,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

		// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*ListCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &ListCommand{}, subCmd)
}

// dylt misc
func TestMainSubMisc(t *testing.T) {
	// Independent values
	cmdName := "dylt"
	cmdFlags := []string{}
	subCmdName := "misc"
	subCmdFlags := []string{}
	subCmdArgs := []string{}

	// Create dependent values for command + test
	cmdArgs := slices.Concat([]string{subCmdName}, subCmdFlags, subCmdArgs)
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t,
		NewMainCommand,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

		// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*MiscCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &MiscCommand{}, subCmd)
}

// dylt misc lookup hostname
func TestMainSubMiscLookupHostname(t *testing.T) {
	// Create + parse main command with "dylt misc lookup hostname"
	subCmdName := "misc"
	subCmdFlags := []string{}
	subCmdArgs := []string{"lookup", "hostname"}
	cmdline := append([]string{"dylt", subCmdName}, subCmdArgs...)
	cmdString := strings.Join(cmdline[0:2], " ")
	cmd := NewMainCommand(cmdline, nil)
	err := cmd.Parse()
	// subcommand & subArgs
	_TestSubCommandAndArgs(t, cmd, subCmdName, subCmdArgs)
	require.NoError(t, err)
	_TestSubcommandCreation[*MiscCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		cmdString,
	)
}

func TestMainSubStatus(t *testing.T) {
	// Independent values
	cmdName := "dylt"
	cmdFlags := []string{}
	subCmdName := "status"
	subCmdFlags := []string{}
	subCmdArgs := []string{}

	// Create dependent values for command + test
	cmdArgs := slices.Concat([]string{subCmdName}, subCmdFlags, subCmdArgs)
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t,
		NewMainCommand,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

		// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*StatusCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &StatusCommand{}, subCmd)
}

// dylt vm
func TestMainSubVm(t *testing.T) {
	// Independent values
	cmdName := "dylt"
	cmdFlags := []string{}
	subCmdName := "vm"
	subCmdFlags := []string{}
	subCmdArgs := []string{}

	// Create dependent values for command + test
	cmdArgs := slices.Concat([]string{subCmdName}, subCmdFlags, subCmdArgs)
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t,
		NewMainCommand,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

	// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*VmCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &VmCommand{}, subCmd)
}

// dylt watch
func TestMainSubWatch(t *testing.T) {
	// Independent values
	cmdName := "dylt"
	cmdFlags := []string{}
	subCmdName := "watch"
	subCmdFlags := []string{}
	subCmdArgs := []string{}

	// Create dependent values for command + test
	cmdArgs := slices.Concat([]string{subCmdName}, subCmdFlags, subCmdArgs)
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t,
		NewMainCommand,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

		// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*WatchCommand](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &WatchCommand{}, subCmd)
}
