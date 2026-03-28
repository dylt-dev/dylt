package cmd

type Command interface {
	Args() (Cmdline, bool)
	CommandLine() Cmdline
	CommandName() string
	CommandArgs() ([]string, bool)
	CommandString() (string, bool)
	CreateSubCommand() (Command, error)
	HandleArgs() error
	Parse() error
	PrintUsage()
	Run() error
	SubArgs() (Cmdline, bool)
	SubCommand() (string, bool)
}
