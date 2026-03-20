package cmd

type Cmdline []string

func (o Cmdline) Args () Cmdline {
	if len(o) < 1 {
		return nil
	}
	var args Cmdline = (o)[1:]
	return args
}

func (o Cmdline) Command () string {
	if len(o) <= 0 {
		return ""
	}
	return (o)[0]
}

func (o Cmdline) HasCommand () bool {
	return len(o) > 0
}

type Command interface {
	Args () (Cmdline, bool)
	CommandName () string
	GetCommandArgs () ([]string, bool)
	GetCommandString () (string, bool)
	HandleArgs () error
	Parse () error
	Run () error
	SubArgs() (Cmdline, bool)
	SubCommand() (string, bool)
}

type SuperCommand interface {
	Command
	CreateSubCommand () (Command, error)
}