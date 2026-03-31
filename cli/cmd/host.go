package cmd

import (
	"fmt"
	"log/slog"

	"github.com/dylt-dev/dylt/api"
	"github.com/dylt-dev/dylt/common"
)

type HostCommand struct {
	*BaseCommand
}

func NewHostCommand(cmdline Cmdline, parent Command) *HostCommand {
	// host command
	name := "host"
	cmd := &HostCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Host)}
	
	//init flags (if any)
	
	return cmd
}

func (cmd *HostCommand) CreateSubCommand() (Command, error) {
	args, is := cmd.Args()
	if !is {
		return nil, nil
	}
	return createHostSubCommand(args, cmd)
}

func (cmd *HostCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}

	// if Help flag is set, no further processing is necessary
	if cmd.Help {
		return nil
	}

	// validate arg count
	cmdArgs, _ := cmd.Args()
	nExpected := 0
	if len(cmdArgs) < nExpected {
		cmd.PrintUsage()
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("`%s` expects >=%d argument(s); received %d",
			cmdString,
			nExpected,
			len(cmdArgs))
	}

	// init positional params (nop - no positional params)

	return nil
}

func (cmd *HostCommand) Run() error {
	slog.Debug("HostCommand.Run()", "args", cmd.Cmdline)

	// Check for 0 args; if so print usage & return
	if len(cmd.Cmdline) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// Parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}

	// If help flag set, print usage + return
	if cmd.Help {
		cmd.PrintUsage()
		return nil
	}

	// If no args, print usage
	args, _ := cmd.Args()
	if len(args) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// Execute command
	err = RunHost(args, cmd)
	return err
}

func RunHost(cmdline Cmdline, parent Command) error {
	slog.Debug("RunHost()", "cmdline", cmdline, "parent", parent)
	// Create the subcommand and run it
	subCmd, err := createHostSubCommand(cmdline, parent)
	if err != nil {
		return err
	}
	err = subCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func createHostSubCommand(cmdline Cmdline, parent Command) (Command, error) {
	cmdName := cmdline.Command()
	cmdMap := CommandMap{
		"init": HostInitCommandF.New,
	}

	cmdFactoryFunc, ok := cmdMap[cmdName]
	if !ok {
		parent.PrintUsage()
		return nil, fmt.Errorf("unrecognized command: %s", cmdName)
	}
		
	cmd := cmdFactoryFunc(cmdline, parent)
	return cmd, nil
}

type HostInitCommand struct {
	*BaseCommand
	Gid int
	Uid int
}

func NewHostInitCommand(cmdline Cmdline, parent Command) *HostInitCommand {
	// host init command
	name := "host.init"
	cmd := &HostInitCommand{BaseCommand: NewBaseCommand(name, cmdline, parent, USG_Host_Init)}
	
	//init flags (if any)
	cmd.IntVar(&cmd.Gid, "gid", 2000, "gid")
	cmd.IntVar(&cmd.Uid, "uid", 2000, "uid")

	return cmd
}

func (cmd HostInitCommand) HandleArgs() error {
	// parse flags
	err := cmd.Parse()
	if err != nil {
		return err
	}

	// if Help flag is set, no further processing is necessary
	if cmd.Help {
		return nil
	}

	// validate arg count
	cmdArgs, _ := cmd.Args()
	nExpected := 0
	if len(cmdArgs) != nExpected {
		cmd.PrintUsage()
		cmdString, _ := cmd.CommandString()
		return fmt.Errorf("%s` expects %d argument(s); received %d",
			cmdString,
			nExpected,
			len(cmdArgs))
	}

	// init positional params (nop - no positional params)

	return nil
}

func (cmd *HostInitCommand) Run() error {
	slog.Debug("HostInitCommand.Run()", "args", cmd.Cmdline)

	// Parse flags & get positional args
	err := cmd.HandleArgs()
	if err != nil {
		return err
	}

	// If help flag set, print usage + return
	if cmd.Help {
		cmd.PrintUsage()
		return nil
	}

	// If no args, print usage
	args, _ := cmd.Args()
	if len(args) == 0 {
		common.Logger.Comment("no args; printing usage")
		cmd.PrintUsage()
		return nil
	}

	// Execute command
	err = api.RunHostInit(cmd.Uid, cmd.Gid)
	return err
}
