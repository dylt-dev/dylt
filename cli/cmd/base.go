package cmd

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"strings"

	"github.com/dylt-dev/dylt/common"
)

type CommandOpts interface {}
type EmptyOpts struct {}

type BaseCommand[Opts CommandOpts] struct {
	*flag.FlagSet
	Parent          Command
	Cmdline         Cmdline
	Usage           string
	argMap          ArgMap
	help            bool
	opts            Opts
	subCommandMap   CommandMap
	validator       CommandValidator
	fnRun           func(*BaseCommand[Opts]) error
	isUsageOnNoArgs bool
}

// type BaseCommandS BaseCommand[string]
// type BaseCommandSA BaseCommand[[]string]

type BaseCommandConfig[Opts CommandOpts] struct {
	name            string
	opts            Opts
	usage           string
	validator       CommandValidator
	fnRun           func(*BaseCommand[Opts]) error
	isUsageOnNoArgs bool
}

func NewBaseCommand[Opts CommandOpts] (cmdline Cmdline,
	                                    parent Command,
	                                    cfg BaseCommandConfig[Opts],
                                       ) *BaseCommand[Opts] {
	cmd := &BaseCommand[Opts]{
		Cmdline:       cmdline,
		Parent:        parent,
		fnRun:         cfg.fnRun,
		opts:          cfg.opts,
		validator:     cfg.validator,
		FlagSet:       flag.NewFlagSet(cfg.name, flag.ExitOnError),
		Usage:         cfg.usage,
	}
	cmd.FlagSet.BoolVar(&cmd.help, "help", false, "give it to me")

	return cmd
}

func (cmd *BaseCommand[_]) ArgMap() ArgMap{
	return cmd.argMap
}

func (cmd *BaseCommand[_]) Args() (Cmdline, bool) {
	if !cmd.FlagSet.Parsed() {
		return nil, false
	}
	return cmd.FlagSet.Args(), true
}

func (cmd *BaseCommand[_]) CommandArgs() ([]string, bool) {
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

func (cmd *BaseCommand[_]) CommandLine() Cmdline {
	return cmd.Cmdline
}

func (cmd *BaseCommand[_]) CommandName() string {
	return cmd.Cmdline.Command()
}

func (cmd *BaseCommand[_]) CommandMap() CommandMap {
	return cmd.subCommandMap
}

func (cmd *BaseCommand[_]) CommandString() (string, bool) {
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

func (cmd *BaseCommand[_]) CommandValidator() CommandValidator {
	return cmd.validator
}

type NoSubcommandsError struct{}

func (o NoSubcommandsError) Error() string { return "No Subcommands" }

func (cmd *BaseCommand[Opts]) CreateSubCommand() (Command, error) {
	common.Logger.Debugf("%s %v\n", "cmd.CommandMap()", cmd.CommandMap())
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

	subCmd := cmdFactoryFunc.(func (Cmdline, Command) Command)(cmdline, cmd)
	return subCmd, nil
}

func (cmd *BaseCommand[_]) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}

	// if Help flag is set, no further processing is necessary
	if cmd.Help() {
		return nil
	}

	// validate args
	cmdArgs, _ := cmd.Args()
	var v CommandValidator = cmd.CommandValidator()
	if !v.IsValid(cmdArgs) {
		cmdString, _ := cmd.CommandString()
		errmsg := v.ErrorMessage(cmdArgs)
		cmd.PrintUsage()
		return fmt.Errorf("`%s` %s", cmdString, errmsg)
	}

	// init positional params, if any
	if cmd.argMap != nil {
		for i, ptr := range cmd.argMap {
			*ptr = cmdArgs[i]
		}
	}

	return nil
}

func (cmd *BaseCommand[_]) Help() bool {
	return cmd.help
}

func (cmd *BaseCommand[T]) Opts() any {
	return cmd.opts
}

func (cmd *BaseCommand[_]) Parse() error {
	common.Logger.Debug("hiii")
	common.Logger.Debugf("cmd=%v\n", cmd)
	common.Logger.Debugf("cmd.FlagSet=%v\n", cmd.FlagSet)
	err := cmd.FlagSet.Parse(cmd.Cmdline.Args())
	if err != nil {
		return err
	}
	return nil
}

func (cmd *BaseCommand[_]) PrintUsage() {
	fmt.Print(cmd.Usage)
}

func (cmd *BaseCommand[_]) Run() error {
	slog.Debug("CallCommand.Run()", "args", cmd.Cmdline)

	// Parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}
	args, _ := cmd.Args()

	// If help flag set, print usage + return
	if cmd.Help() {
		cmd.PrintUsage()
		return nil
	}

	// Check for 0 args; if so print usage & return
	common.Logger.Debugf("args: %#+v\n", args)
	common.Logger.Debugf("len(args): %#+v\n", len(args))
	common.Logger.Debugf("cmd.UsageOnNoArgs(): %#+v\n", cmd.UsageOnNoArgs())
	if len(args) == 0 && cmd.UsageOnNoArgs() {
		common.Logger.Debug("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// if CommandMap exists run subcommand
	cmdMap := cmd.CommandMap()
	if cmdMap != nil && len(args) > 0 {
		subCmd, err := cmd.CreateSubCommand()
		if err != nil {
			return err
		}
		err = subCmd.Run()
		return err
	}

	// execute command
	fmt.Printf("cmd.fnRun=%v\n", cmd.fnRun)
	if cmd.fnRun != nil {
		return cmd.fnRun(cmd)
	}

	return nil
}

func (cmd *BaseCommand[_]) SubArgs() (Cmdline, bool) {
	if !cmd.FlagSet.Parsed() {
		return nil, false
	}
	var subCmdline Cmdline = cmd.FlagSet.Args()
	return subCmdline.Args(), true
}

func (cmd *BaseCommand[_]) SubCommand() (string, bool) {
	if !cmd.Parsed() {
		return "", false
	}
	var subCmdline Cmdline = cmd.FlagSet.Args()
	return subCmdline.Command(), true
}

func (cmd *BaseCommand[_]) UsageOnNoArgs() bool {
	return cmd.isUsageOnNoArgs
}
