package cli

import (
	"bufio"
	"flag"
	"io"
	"os"
	"strings"
	"testing"

	clicmd "github.com/dylt-dev/dylt/cli/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type SimpleCommand struct {
	*clicmd.BaseCommand[clicmd.EmptyOpts ]
}

func (cmd SimpleCommand) HandleArgs() error { return nil }
func (cmd SimpleCommand) Run() error        { return nil }

func TestCommandArgs(t *testing.T) {
	cmdline := clicmd.Cmdline{"foo", "bar", "bum"}
	cmd := SimpleCommand{BaseCommand: &clicmd.BaseCommand[clicmd.EmptyOpts]{Cmdline: cmdline, FlagSet: &flag.FlagSet{}}}
	err := cmd.Parse()
	args, _ := cmd.Args()
	t.Logf("args=%v", args)
	require.NoError(t, err)
	_TestArgs(t, cmd, len(cmdline.Args()), cmdline.Args())
}

type InitCommand struct {
	*flag.FlagSet
	EtcdDomain string
}

func NewInitCommand() *InitCommand {
	cmd := InitCommand{
		FlagSet: flag.NewFlagSet("init", flag.PanicOnError),
	}
	cmd.FlagSet.StringVar(&cmd.EtcdDomain, "etcd-domain", "", "etcd-domain")
	return &cmd
}

// type arglist []string

// func (o *arglist) Command() string {
// 	return (*o)[0]
// }

// func (o *arglist) Args() []string {
// 	return (*o)[1:]
// }

func TestArgs(t *testing.T) {
	cmdline := clicmd.Cmdline{"foo", "bar", "bum"}
	flagSet := flag.FlagSet{}
	flagSet.Parse(cmdline)
	t.Logf("flagSet.Args(): %#+v\n", flagSet.Args())
	args := flagSet.Args()
	require.Equal(t, len(cmdline), len(args))
	require.Equal(t, cmdline, clicmd.Cmdline(args))
}

func TestArgsFlag(t *testing.T) {
	cmdline := clicmd.Cmdline{"dylt", "init", "--etcd-domain", "foo.dylt.dev"}
	flagSet := flag.FlagSet{}
	var etcdDomain string
	flagSet.StringVar(&etcdDomain, "etcd-domain", "", "")
	err := flagSet.Parse(cmdline[1:])
	require.NoError(t, err)
	t.Logf("flagSet.Args()=%v", flagSet.Args())
}

func TestBoolFlag(t *testing.T) {
	var cmdline clicmd.Cmdline = []string{"--help"}
	var fs flag.FlagSet
	var helpy *bool = fs.Bool("help", false, "yeah yeah yeah")
	fs.Parse(cmdline)
	t.Log(fs.Args())
	require.True(t, *helpy)
}

func TestInitCommand(t *testing.T) {
	cmdline := clicmd.Cmdline(strings.Split("dylt init --etcd-domain hello.dylt.dev", " "))
	mainCmd := flag.NewFlagSet("dylt", flag.PanicOnError)
	mainCmd.Parse(cmdline.Args())
	var args clicmd.Cmdline = mainCmd.Args()
	cmd := args.Command()
	assert.Equal(t, "init", cmd)
	initCmd := clicmd.InitCommandF.New(args, nil).(*clicmd.BaseCommand[clicmd.InitOpts])
	initCmd.Parse()
	opts, is := initCmd.Opts().(*clicmd.InitOpts)
	require.True(t, is)
	require.NotNil(t, opts)
	assert.Equal(t, "hello.dylt.dev", opts.EtcdDomain)

	t.Log(strings.Join(args, " "))
}

func TestLongLine(t *testing.T) {
	args := strings.Split("dylt call --script /opt/bin/daylight.sh github-release-list --token ghuBlahDiBlah dylt-dev yellowrose", " ")
	assert.Equal(t, "dylt", args[0])
	assert.Equal(t, "call", args[1])
	cmdCall := flag.NewFlagSet("call", flag.PanicOnError)
	var scriptPath string
	cmdCall.StringVar(&scriptPath, "script", "", "Path to script")
	cmdCall.Parse(args[2:])
	argsCall := cmdCall.Args()
	assert.Equal(t, "github-release-list", argsCall[0])
	assert.Equal(t, "/opt/bin/daylight.sh", scriptPath)
	cmdGhr := flag.NewFlagSet("github-release-list", flag.PanicOnError)
	var token string
	cmdGhr.StringVar(&token, "token", "", "token")
	cmdGhr.Parse(argsCall[1:])
	argsGhr := cmdGhr.Args()
	assert.Equal(t, "ghuBlahDiBlah", token)
	assert.Equal(t, "dylt-dev", argsGhr[0])
	assert.Equal(t, "yellowrose", argsGhr[1])
}

func TestMissingBool(t *testing.T) {
	cmdline := clicmd.Cmdline{"dylt"}
	flagSet := flag.FlagSet{}
	var help bool
	flagSet.BoolVar(&help, "help", true, "help")
	flagSet.Parse(cmdline.Args())
	t.Logf("flagSet.Args()=%v", flagSet.Args())
	require.True(t, help)
}

func TestParseBool(t *testing.T) {
	cmdline := clicmd.Cmdline{"dylt", "--help"}
	flagSet := flag.FlagSet{}
	var help bool
	flagSet.BoolVar(&help, "help", false, "help")
	flagSet.Parse(cmdline.Args())
	t.Logf("flagSet.Args()=%v", flagSet.Args())
	require.True(t, help)
}

func TestStringFlag(t *testing.T) {
	cmdline := clicmd.Cmdline{"dylt", "--foo", "bar"}
	var fs flag.FlagSet
	var helpy *string = fs.String("foo", "", "yeah yeah yeah")
	fs.Parse(cmdline.Args())
	require.Equal(t, "bar", *helpy)
}

// cmdline: dylt --foo 1 --bar x y
// Expected:
func TestWeirdArgs(t *testing.T) {
	args := strings.Split("dylt --foo 1 --bar x y", " ")
	flagset := flag.NewFlagSet("dylt", flag.PanicOnError)
	var fooptr *string = flagset.String("foo", "", "foo")
	var barptr *string = flagset.String("bar", "", "bar")
	flagset.Parse(args[1:])
	require.Equal(t, "1", *fooptr)
	require.Equal(t, "x", *barptr)
	argz := flagset.Args()
	require.Equal(t, 1, len(argz))
	require.Equal(t, "y", argz[0])
}

func TestRead(t *testing.T) {
	if os.Getenv("INTERACTIVE") != "Y" {
		t.Skip("Test run is non-interactive; skipping test")
	}
	stdin, _ := io.Pipe()
	scanner := bufio.NewScanner(stdin)
	scanner.Scan()
	s := scanner.Text()
	err := scanner.Err()
	assert.NoError(t, err)
	assert.NotEmpty(t, s)
}

func TestSubcommand(t *testing.T) {
	cmdline := []string{"foo", "bar", "bum"}
	flagSet := flag.FlagSet{}
	flagSet.Parse(cmdline)
	t.Logf("flagSet.Args()=%v", flagSet.Args())
}

type PappyCommand struct{ *clicmd.BaseCommand[clicmd.EmptyOpts] }
type DaddyCommand struct{ *clicmd.BaseCommand[clicmd.EmptyOpts] }
type MeCommand struct{ *clicmd.BaseCommand[clicmd.EmptyOpts] }

func (cmd *PappyCommand) HandleArgs() error { return nil }
func (cmd *PappyCommand) Run() error        { return nil }

type EC clicmd.Command

func NewPappyCommand(cmdline clicmd.Cmdline, parent clicmd.Command) *PappyCommand {
	return &PappyCommand{BaseCommand: &clicmd.BaseCommand[clicmd.EmptyOpts]{Cmdline: cmdline, FlagSet: &flag.FlagSet{}}}
}

func (cmd *DaddyCommand) HandleArgs() error { return nil }
func (cmd *DaddyCommand) Run() error        { return nil }

func NewDaddyCommand(cmdline clicmd.Cmdline, parent clicmd.Command) *DaddyCommand {
	return &DaddyCommand{BaseCommand: &clicmd.BaseCommand[clicmd.EmptyOpts]{Cmdline: cmdline, FlagSet: &flag.FlagSet{}, Parent: parent}}
}

func (cmd *MeCommand) HandleArgs() error { return nil }
func (cmd *MeCommand) Run() error        { return nil }

func NewMeCommand(cmdline clicmd.Cmdline, parent clicmd.Command) *MeCommand {
	return &MeCommand{BaseCommand: &clicmd.BaseCommand[clicmd.EmptyOpts]{Cmdline: cmdline, FlagSet: &flag.FlagSet{}, Parent: parent}}
}

type ECF clicmd.CommandFactory
var PappyCommandF ECF = ECF{ FnNew: func (cmdline clicmd.Cmdline, parent clicmd.Command) clicmd.Command { return NewPappyCommand(cmdline, parent) }}
var DaddyCommandF ECF = ECF{FnNew: func (cmdline clicmd.Cmdline, parent clicmd.Command) clicmd.Command { return NewDaddyCommand(cmdline, parent) }}
var MeCommandF ECF = ECF{FnNew: func (cmdline clicmd.Cmdline, parent clicmd.Command) clicmd.Command { return NewMeCommand(cmdline, parent) }}

func TestPappyDaddyMe(t *testing.T) {
	var cmdline clicmd.Cmdline = []string{"pappy", "daddy", "me", "foo"}

	// Create `pappy` subcommand
	pappy := PappyCommandF.FnNew(cmdline, nil)
	_TestPreParseValues(t, pappy)
	// Parse and check again
	pappy.Parse()
	_TestSubCommandAndArgs(t, pappy, "daddy", []string{"me", "foo"})
	_TestCommandString(t, pappy, "pappy")

	// Create `daddy` subcommand
	daddyArgs, is := pappy.Args()
	require.True(t, is)
	daddy := DaddyCommandF.FnNew(daddyArgs, pappy)
	_TestPreParseValues(t, daddy)
	// Parse and check again
	daddy.Parse()
	_TestSubCommandAndArgs(t, daddy, "me", []string{"foo"})
	_TestCommandString(t, daddy, "pappy daddy")

	// Create `me` subcommand
	meArgs, is := daddy.Args()
	me := MeCommandF.FnNew(meArgs, daddy)
	_TestPreParseValues(t, me)
	// Parse and check again
	me.Parse()
	_TestSubCommandAndArgs(t, me, "foo", []string{})
	_TestCommandString(t, me, "pappy daddy me")
}

func _TestArgs(t *testing.T, cmd clicmd.Command, targetArgsLen int, targetArgs clicmd.Cmdline) {
	args, is := cmd.Args()
	require.True(t, is)
	require.Equal(t, targetArgsLen, len(args))
	require.Equal(t, targetArgs, args)
}

func _TestCommandString(t *testing.T, cmd clicmd.Command, targetCmdString string) {
	cmdString, flag := cmd.CommandString()
	require.True(t, flag)
	require.Equal(t, targetCmdString, cmdString)
}

func _TestPreParseValues(t *testing.T, cmd clicmd.Command) {
	// Test pre-parse values are as expected
	subCommand, flag := cmd.SubCommand()
	require.Empty(t, subCommand)
	require.False(t, flag)
	subArgs, flag := cmd.SubArgs()
	require.Nil(t, subArgs)
	require.False(t, flag)
}

func _TestSubCommandAndArgs(t *testing.T, cmd clicmd.Command, targetSubCommand string, targetSubArgs clicmd.Cmdline) {
	subCommand, flag := cmd.SubCommand()
	require.Equal(t, targetSubCommand, subCommand)
	require.True(t, flag)
	subArgs, flag := cmd.SubArgs()
	require.Equal(t, targetSubArgs, subArgs)
	require.True(t, flag)
}
