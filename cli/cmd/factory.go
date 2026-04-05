package cmd


// Simple method signature for creating a command
type NewCommandFunc[Opts CommandOpts] func (Cmdline, Command[Opts]) Command[Opts]

// Simple interface for Command Factory objects, that create a single
// type of Command (eg CallCommand)
//
// U - the type of Command created by this factory
type ICommandFactory[Opts CommandOpts] interface {
	New (Cmdline, Command[Opts])
}

// Simple implementation of ICommandFactory. It has a single field - a 
// reference to a function that returns a type of Command.
// This allows the creation of CommandFactories that are parameterized
// on a Command type and on a factory function for creating that command type.
type CommandFactory[Opts CommandOpts] struct {
	FnNew NewCommandFunc[Opts]
}

func (cf CommandFactory[Opts]) New (cmdline Cmdline, parent Command[Opts]) Command[Opts] {
	return cf.FnNew(cmdline, parent)
}

// Factories for each type of command

var CallCommandF CommandFactory[CallOpts] = CommandFactory[CallOpts]{ FnNew: func (cmdline Cmdline, parent Command[CallOpts]) Command[CallOpts] { return NewCallCommand(cmdline, parent) }}
var ConfigCommandF CommandFactory[ConfigOpts] = CommandFactory[ConfigOpts]{ FnNew: func (cmdline Cmdline, parent Command[ConfigOpts]) Command[ConfigOpts] { return NewConfigCommand(cmdline, parent) } }
var ConfigGetCommandF CommandFactory[ConfigGetOpts] = CommandFactory[ConfigGetOpts]{ FnNew: func (cmdline Cmdline, parent Command[ConfigGetOpts]) Command[ConfigGetOpts] { return NewConfigGetCommand(cmdline, parent) } }
var ConfigSetCommandF CommandFactory[ConfigSetOpts] = CommandFactory[ConfigSetOpts]{ FnNew: func (cmdline Cmdline, parent Command[ConfigSetOpts]) Command[ConfigSetOpts] { return NewConfigSetCommand(cmdline, parent) } }
var ConfigShowCommandF CommandFactory[ConfigShowOpts] = CommandFactory[ConfigShowOpts]{ FnNew: func (cmdline Cmdline, parent Command[ConfigShowOpts]) Command[ConfigShowOpts] { return NewConfigShowCommand(cmdline, parent) } }
var CreateTwoNodeClusterCommandF CommandFactory = CommandFactory{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewCreateTwoNodeClusterCommand(cmdline, parent) } }
var GenEtcdRunScriptCommandF CommandFactory = CommandFactory{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewGenEtcdRunScriptCommand(cmdline, parent) } }
var GetCommandF CommandFactory[GetOpts] = CommandFactory[GetOpts]{ FnNew: func (cmdline Cmdline, parent Command[GetOpts]) Command[GetOpts] { return NewGetCommand(cmdline, parent) }}
var HostCommandF CommandFactory[HostOpts] = CommandFactory[HostOpts]{ FnNew: func (cmdline Cmdline, parent Command[HostOpts]) Command[HostOpts] { return NewHostCommand(cmdline, parent) } }
var HostInitCommandF CommandFactory[HostInitOpts] = CommandFactory[HostInitOpts]{ FnNew: func (cmdline Cmdline, parent Command[HostInitOpts]) Command[HostInitOpts] { return NewHostInitCommand(cmdline, parent) } }
var InitCommandF CommandFactory = CommandFactory{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewInitCommand(cmdline, parent) }}
var ListCommandF CommandFactory = CommandFactory{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewListCommand(cmdline, parent) }}
var LookupCommandF CommandFactory = CommandFactory{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewLookupCommand(cmdline, parent) } }
var MainCommandF CommandFactory = CommandFactory{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewMainCommand(cmdline, parent) } }
var MiscCommandF CommandFactory = CommandFactory{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewMiscCommand(cmdline, parent) } }
var StatusCommandF CommandFactory = CommandFactory{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewStatusCommand(cmdline, parent) } }
var VmAddCommandF CommandFactory = CommandFactory{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewVmAddCommand(cmdline, parent) } }
var VmAllCommandF CommandFactory = CommandFactory{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewVmAllCommand(cmdline, parent) } }
var VmCommandF CommandFactory = CommandFactory{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewVmCommand(cmdline, parent) } }
var VmDelCommandF CommandFactory = CommandFactory{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewVmDelCommand(cmdline, parent) } }
var VmGetCommandF CommandFactory = CommandFactory{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewVmGetCommand(cmdline, parent) } }
var VmListCommandF CommandFactory = CommandFactory{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewVmListCommand(cmdline, parent) } }
var VmSetCommandF CommandFactory = CommandFactory{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewVmSetCommand(cmdline, parent) } }
var WatchCommandF CommandFactory = CommandFactory{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewWatchCommand(cmdline, parent) } }
var WatchScriptCommandF CommandFactory = CommandFactory{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewWatchScriptCommand(cmdline, parent) } }
var WatchSvcCommandF CommandFactory = CommandFactory{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewWatchSvcCommand(cmdline, parent) } }
// var CallCommandF CommandFactory[*CallCommand] = CommandFactory[*CallCommand]{ FnNew: NewCallCommand }
// var CallCommandF CommandFactory[*CallCommand] = CommandFactory[*CallCommand]{ FnNew: NewCallCommand }
// var CallCommandF CommandFactory[*CallCommand] = CommandFactory[*CallCommand]{ FnNew: NewCallCommand }
// var CallCommandF CommandFactory[*CallCommand] = CommandFactory[*CallCommand]{ FnNew: NewCallCommand }
// var CallCommandF CommandFactory[*CallCommand] = CommandFactory[*CallCommand]{ FnNew: NewCallCommand }
// var CallCommandF CommandFactory[*CallCommand] = CommandFactory[*CallCommand]{ FnNew: NewCallCommand }
// var CallCommandF CommandFactory[*CallCommand] = CommandFactory[*CallCommand]{ FnNew: NewCallCommand }
// var CallCommandF CommandFactory[*CallCommand] = CommandFactory[*CallCommand]{ FnNew: NewCallCommand }
// var CallCommandF CommandFactory[*CallCommand] = CommandFactory[*CallCommand]{ FnNew: NewCallCommand }
// var CallCommandF CommandFactory[*CallCommand] = CommandFactory[*CallCommand]{ FnNew: NewCallCommand }
// var CallCommandF CommandFactory[*CallCommand] = CommandFactory[*CallCommand]{ FnNew: NewCallCommand }
