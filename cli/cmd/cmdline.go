package cmd

type Cmdline []string

func NewCmdline(name string, flags []string, args []string) Cmdline {
	var cmdline Cmdline
	cmdline = Cmdline{name}
	cmdline = append(cmdline, flags...)
	cmdline = append(cmdline, args...)
	return cmdline
}

func (o Cmdline) Args() Cmdline {
	if len(o) < 1 {
		return nil
	}
	var args Cmdline = (o)[1:]
	return args
}

func (o Cmdline) Command() string {
	if len(o) <= 0 {
		return ""
	}
	return (o)[0]
}

func (o Cmdline) HasCommand() bool {
	return len(o) > 0
}

type Command interface {
	Args() (Cmdline, bool)
	CommandLine() Cmdline
	CommandName() string
	CommandArgs() ([]string, bool)
	CommandString() (string, bool)
	HandleArgs() error
	Parse() error
	Run() error
	SubArgs() (Cmdline, bool)
	SubCommand() (string, bool)
}

type SuperCommand interface {
	Command
	CreateSubCommand() (Command, error)
}
