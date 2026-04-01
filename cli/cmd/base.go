package cmd

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type BaseCommand struct {
	*flag.FlagSet
	Parent  Command
	Cmdline Cmdline
	Help    bool
	Usage string
	argmap map[int]*string
	commandMap CommandMap
	commandValidator CommandValidator
	fnRun func () error
	isUsageOnNoArgs bool
}

// type BaseCommandS BaseCommand[string]
// type BaseCommandSA BaseCommand[[]string]

func NewBaseCommand[U UsageTextType](name string,
                                     cmdline Cmdline,
									 parent Command,
									 usageText U,
									 cmdMap CommandMap,
									 cmdValidator CommandValidator,
								    ) *BaseCommand {
	cmd := &BaseCommand{
		Cmdline: cmdline,
		commandMap: cmdMap,
		commandValidator: cmdValidator,
		Parent:  parent,
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),
		Usage: CreateUsageString(usageText),
	}
	cmd.FlagSet.BoolVar(&cmd.Help, "help", false, "give it to me")

	return cmd
}

func (cmd *BaseCommand) Args() (Cmdline, bool) {
	if !cmd.FlagSet.Parsed() {
		return nil, false
	}
	return cmd.FlagSet.Args(), true
}

func (cmd *BaseCommand) CommandArgs() ([]string, bool) {
	if !cmd.FlagSet.Parsed() {
		return nil, false
	}
	var cmdArgs []string
	var flag bool
	if cmd.Parent == nil {
		cmdArgs = []string{}
	} else {
		cmdArgs, flag = cmd.Parent.CommandArgs()
		if !flag {
			return nil, flag
		}
	}
	cmdArgs = append(cmdArgs, cmd.Cmdline.Command())
	return cmdArgs, true
}

func (cmd *BaseCommand) CommandLine() Cmdline {
	return cmd.Cmdline
}

func (cmd *BaseCommand) CommandName() string {
	return cmd.Cmdline.Command()
}

func (cmd *BaseCommand) CommandMap() CommandMap {
	return cmd.commandMap
}

func (cmd *BaseCommand) CommandString() (string, bool) {
	if !cmd.FlagSet.Parsed() {
		return "", false
	}
	cmdArgs, flag := cmd.CommandArgs()
	if !flag {
		return "", false
	}
	if cmdArgs == nil {
		return "", flag
	}
	return strings.Join(cmdArgs, " "), true
}

func (cmd *BaseCommand) CommandValidator() CommandValidator {
	return cmd.commandValidator
}

type NoSubcommandsError struct {}
func (o NoSubcommandsError) Error() string { return "No Subcommands" }

func (cmd *BaseCommand) CreateSubCommand () (Command, error) {
	fmt.Printf("%s %v\n", "cmd.CommandMap()", cmd.CommandMap())
	cmdline, is := cmd.Args()
	if !is {
		return nil, errors.New("Command not Parse()'d")
	}
	if cmd.CommandMap() == nil {
		return nil, nil
	}
	// return createConfigSubCommand(args, cmd)
	cmdName := cmdline.Command()
	cmdMap := cmd.CommandMap()
	if cmdMap == nil {
		return nil, nil
	}

	cmdFactoryFunc, ok := cmdMap[cmdName]
	if !ok {
		cmd.PrintUsage()
		return nil, fmt.Errorf("unrecognized command: %s", cmdName)
	}
		
	subCmd := cmdFactoryFunc(cmdline, cmd)
	return subCmd, nil
}

func (cmd *BaseCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}

	// if Help flag is set, no further processing is necessary
	if cmd.Help {
		return nil
	}

	// validate args
	cmdArgs, _ := cmd.Args()
	var v CommandValidator = cmd.CommandValidator()
	if ! v.IsValid(cmdArgs) {
		cmdString, _ := cmd.CommandString()
		errmsg := v.ErrorMessage(cmdArgs)
		return fmt.Errorf("`%s` %s", cmdString, errmsg)
	}

	// init positional params, if any
	if cmd.argmap != nil {
		for i, ptr := range cmd.argmap {
			*ptr = cmdArgs[i]
		}
	}

	return nil
}


func (cmd *BaseCommand) Parse() error {
	fmt.Printf("cmd.FlagSet=%v\n", cmd.FlagSet)
	err := cmd.FlagSet.Parse(cmd.Cmdline.Args())
	if err != nil {
		return err
	}
	return nil
}

func (cmd *BaseCommand) PrintUsage () {
	fmt.Print(cmd.Usage)
	fmt.Println()
}

func (cmd *BaseCommand) Run() error { return nil }

func (cmd *BaseCommand) SubArgs() (Cmdline, bool) {
	if !cmd.FlagSet.Parsed() {
		return nil, false
	}
	var subCmdline Cmdline = cmd.FlagSet.Args()
	return subCmdline.Args(), true
}

func (cmd *BaseCommand) SubCommand() (string, bool) {
	if !cmd.Parsed() {
		return "", false
	}
	var subCmdline Cmdline = cmd.FlagSet.Args()
	return subCmdline.Command(), true
}

func (cmd *BaseCommand) UsageOnNoArgs () bool {
	return cmd.isUsageOnNoArgs
}
