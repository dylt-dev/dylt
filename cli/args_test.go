package cli

import (
	"bufio"
	"flag"
	"io"
	"strings"
	"testing"

	clicmd "github.com/dylt-dev/dylt/cli/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

type arglist []string

func (o *arglist) Command() string {
	return (*o)[0]
}

func (o *arglist) Args() []string {
	return (*o)[1:]
}

func TestBoolFlag(t *testing.T) {
	var cmdline clicmd.Cmdline = []string{"dylt", "--help"}
	var fs flag.FlagSet
	var helpy *bool = fs.Bool("help", false, "yeah yeah yeah")
	t.Log(cmdline.Args())
	fs.Parse(cmdline.Args())
	require.True(t, *helpy)
}

func TestInitCommand(t *testing.T) {
	cmdline := arglist(strings.Split("dylt init --etcd-domain hello.dylt.dev", " "))
	mainCmd := flag.NewFlagSet("dylt", flag.PanicOnError)
	mainCmd.Parse(cmdline.Args())
	var args arglist = mainCmd.Args()
	cmd := args.Command()
	assert.Equal(t, "init", cmd)
	initCmd := NewInitCommand()
	initCmd.Parse(args.Args())
	assert.Equal(t, "hello.dylt.dev", initCmd.EtcdDomain)

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

type PappyCommand struct { *clicmd.BaseCommand }
type DaddyCommand struct{ *clicmd.BaseCommand }
type MeCommand struct{ *clicmd.BaseCommand }

func (cmd *PappyCommand) Run () error { return nil}

func NewPappyCommand (cmdline clicmd.Cmdline) *PappyCommand {
	return &PappyCommand{BaseCommand: &clicmd.BaseCommand{Cmdline: cmdline, FlagSet: &flag.FlagSet{}}}
}

func (cmd *DaddyCommand) Run () error { return nil }

func NewDaddyCommand (cmdline clicmd.Cmdline, parent *PappyCommand) *DaddyCommand {
	return &DaddyCommand{BaseCommand: &clicmd.BaseCommand{Cmdline: cmdline, FlagSet: &flag.FlagSet{}, ParentCommand: parent}}
}

func (cmd *MeCommand) Run () error { return nil }

func NewMeCommand (cmdline clicmd.Cmdline, parent *DaddyCommand) *MeCommand {
	return &MeCommand{BaseCommand: &clicmd.BaseCommand{Cmdline: cmdline, FlagSet: &flag.FlagSet{}, ParentCommand: parent}}
}


func TestPappyDaddyMe(t *testing.T) {
	var subCommand string
	var subArgs clicmd.Cmdline
	var flag bool
	var cmdline clicmd.Cmdline = []string{"pappy", "daddy", "me", "foo"}
	pappy := NewPappyCommand(cmdline)
	// Test pappy values pre-parse are as expected
	subCommand, flag = pappy.SubCommand()
	require.Empty(t, subCommand)
	require.False(t, flag)
	subArgs, flag = pappy.SubArgs()
	require.Nil(t, subArgs)
	require.False(t, flag)
	// Parse and check again
	pappy.Parse()
	subCommand, flag = pappy.SubCommand()
	require.Equal(t, "daddy", subCommand)
	require.True(t, flag)
	subArgs, flag = pappy.SubArgs()
	require.Equal(t, clicmd.Cmdline{"me", "foo"}, subArgs)	
	require.True(t, flag)
	// Create `daddy` subcommand
	daddy := NewDaddyCommand(pappy.Cmdline.Args(), pappy)
	// Test daddy values pre-parse are as expected
	subCommand, flag = daddy.SubCommand()
	require.Empty(t, subCommand)
	require.False(t, flag)
	subArgs, flag = daddy.SubArgs()
	require.Nil(t, subArgs)
	require.False(t, flag)
	// Parse and check again
	daddy.Parse()
	subCommand, flag = daddy.SubCommand()
	require.Equal(t, "me", subCommand)
	require.True(t, flag)
	subArgs, flag = daddy.SubArgs()
	require.Equal(t, clicmd.Cmdline{"foo"}, subArgs)	
	require.True(t, flag)
	// Create `me` subcommand
	me := NewMeCommand(daddy.Cmdline.Args(), daddy)
	// Test daddy values pre-parse are as expected
	subCommand, flag = me.SubCommand()
	require.Empty(t, subCommand)
	require.False(t, flag)
	subArgs, flag = me.SubArgs()
	require.Nil(t, subArgs)
	require.False(t, flag)
	// Parse and check again
	me.Parse()
	args, flag := me.Args()
	require.Equal(t, 1, len(args.Args()))
	require.Equal(t, "foo", args.Args()[0])
	require.True(t, flag)
	subArgs, flag = me.SubArgs()
	require.Empty(t, subArgs)
	require.True(t, flag)
	
}