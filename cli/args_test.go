package cli

import (
	"flag"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type InitCommand struct {
	*flag.FlagSet
	EtcdDomain string
}

func NewInitCommand () *InitCommand {
	cmd := InitCommand {
		FlagSet: flag.NewFlagSet("init", flag.PanicOnError),
	}
	cmd.FlagSet.StringVar(&cmd.EtcdDomain, "etcd-domain", "", "etcd-domain")
	return &cmd
}

type arglist []string

func (o *arglist) Command () string {
	return (*o)[0]
}

func (o *arglist) Args () []string {
	return (*o)[1:]
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


func TestLongLine (t *testing.T) {
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

func TestWeirdArgs (t *testing.T) {
	args := strings.Split("dylt --foo 1 --bar x y", " ")
	dylt := flag.NewFlagSet("dylt", flag.PanicOnError)
	dylt.String("foo", "", "foo")
	dylt.String("bar", "", "bar")
	dylt.Parse(args[1:])
	argz := dylt.Args()
	t.Log(strings.Join(argz, " "))
}