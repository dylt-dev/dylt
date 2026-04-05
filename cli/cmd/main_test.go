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
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdName := "dylt"
	cmdFlags := []string{}
	cmdArgs := []string{}
	cmdString := CreateCommandString(cmdName, cmdArgs)
	cmd, is := CreateAndTestCommand(t, MainCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString).(*BaseCommand[MainOpts])
	require.True(t, is)
	require.NotNil(t, cmd)
	require.False(t, cmd.Help())
}


// dylt --help
func TestHelp (t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
	cmdName := "dylt"
	cmdFlags := []string{"--help"}
	cmdArgs := []string{}
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd, is := CreateAndTestCommand(t, MainCommandF.New, cmdName, cmdFlags, cmdArgs, cmdString).(*BaseCommand[MainOpts])
	require.True(t, is)
	require.NotNil(t, cmd)
	require.True(t, cmd.Help())
}


// dylt call foo
func TestMainSubCall(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
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
		MainCommandF.New,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

		// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*BaseCommand[CallOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &BaseCommand[CallOpts]{}, subCmd)
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
		MainCommandF.New,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

		// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*BaseCommand[ConfigOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &BaseCommand[ConfigOpts]{}, subCmd)
}


// dylt config get foo
func TestMainSubConfigGetFoo(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
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
		MainCommandF.New,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

		// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*BaseCommand[ConfigOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &BaseCommand[ConfigOpts]{}, subCmd)
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
		MainCommandF.New,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

		// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*BaseCommand[GetOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &BaseCommand[GetOpts]{}, subCmd)
}


func TestMainSubHost(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
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
		MainCommandF.New,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

	// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*BaseCommand[HostOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &BaseCommand[HostOpts]{}, subCmd)
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
		MainCommandF.New,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

	// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*BaseCommand[HostOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &BaseCommand[HostOpts]{}, subCmd)
}


// dylt init --etcd-domain foo.dylt.dev
func TestMainSubInit(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
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
		MainCommandF.New,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

	// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*BaseCommand[InitOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.Equal(t, etcdDomain, subCmd.opts.EtcdDomain)
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
		MainCommandF.New,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

		// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*BaseCommand[ListOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &BaseCommand[ListOpts]{}, subCmd)
}

// dylt misc
func TestMainSubMisc(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
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
		MainCommandF.New,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

		// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*BaseCommand[MiscOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &BaseCommand[MiscOpts]{}, subCmd)
}

// dylt misc lookup hostname
func TestMainSubMiscLookupHostname(t *testing.T) {
	// Independent values
	cmdName := "dylt"
	cmdFlags := []string{}
	subCmdName := "misc"
	subCmdFlags := []string{}
	subCmdArgs := []string{"lookup", "hostname"}

	// Create dependent values for command + test
	cmdArgs := slices.Concat([]string{subCmdName}, subCmdFlags, subCmdArgs)
	cmdString := fmt.Sprintf("%s", cmdName)
	cmd := CreateAndTestCommand(t,
		MainCommandF.New,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

		// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*BaseCommand[MiscOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &BaseCommand[MiscOpts]{}, subCmd)
}

func TestMainSubStatus(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
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
		MainCommandF.New,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

		// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*BaseCommand[StatusOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &BaseCommand[StatusOpts]{}, subCmd)
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
		MainCommandF.New,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

	// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*BaseCommand[VmOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &BaseCommand[VmOpts]{}, subCmd)
}

// dylt watch
func TestMainSubWatch(t *testing.T) {
	fnTeardown := setup(t)
	defer fnTeardown(t)
	
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
		MainCommandF.New,
		cmdName,
		cmdFlags,
		cmdArgs,
		cmdString)

		// Create dependent values for subcommand + test
	subCmdString := fmt.Sprintf("%s %s", cmdName, subCmdName)
	subCmd := _TestSubcommandCreation[*BaseCommand[WatchOpts]](t,
		cmd,
		subCmdName,
		subCmdFlags,
		subCmdArgs,
		subCmdString,
	)
	require.IsType(t, &BaseCommand[WatchOpts]{}, subCmd)
}
