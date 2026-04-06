package cmd

import "fmt"

type Command interface {
	ArgMap() (ArgMap)
	Args() (Cmdline, bool)
	CommandLine() Cmdline
	CommandName() string
	CommandArgs() ([]string, bool)
	CommandString() (string, bool)
	CommandMap() CommandMap
	CommandValidator() CommandValidator
	CreateSubCommand() (Command, error)
	HandleArgs() error
	Help() bool
	Opts() any
	Parse() error
	PrintUsage()
	Run() error
	SubArgs() (Cmdline, bool)
	SubCommand() (string, bool)
	UsageOnNoArgs() bool
}

type ArgMap map[int]*string
type CommandMap map[string]NewCommandFunc

type CommandValidator interface {
	IsValid(args []string) bool
	ErrorMessage(args []string) string
}

type ArgCountValidator struct {
	CommandValidator
	nExpected int
}

func (v ArgCountValidator) IsValid(args []string) bool {
	// validate arg count
	nArgs := len(args)
	return nArgs == v.nExpected
}

func (v ArgCountValidator) ErrorMessage(args []string) string {
	if v.IsValid(args) {
		return ""
	}

	nArgs := len(args)
	return fmt.Sprintf("expects %d argument(s); received %d",
		v.nExpected, nArgs)
}

type ArgCountGEValidator struct {
	CommandValidator
	nExpected int
}

func (v ArgCountGEValidator) IsValid(args []string) bool {
	// validate arg count
	nArgs := len(args)
	return nArgs >= v.nExpected
}

func (v ArgCountGEValidator) ErrorMessage(args []string) string {
	if v.IsValid(args) {
		return ""
	}

	nArgs := len(args)
	return fmt.Sprintf("expects >= %d argument(s); received %d",
		v.nExpected, nArgs)
}
