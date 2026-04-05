package cmd


// Simple method signature for creating a command
type NewCommandFunc[Opts CommandOpts] func (Cmdline, Command) Command

// Simple interface for Command Factory objects, that create a single
// type of Command (eg CallCommand)
//
// U - the type of Command created by this factory
type ICommandFactory[Opts CommandOpts] interface {
	New (Cmdline, Command)
}

// Simple implementation of ICommandFactory. It has a single field - a 
// reference to a function that returns a type of Command.
// This allows the creation of CommandFactories that are parameterized
// on a Command type and on a factory function for creating that command type.
type CommandFactory[Opts CommandOpts] struct {
	FnNew NewCommandFunc[Opts]
}

func (cf CommandFactory[Opts]) New (cmdline Cmdline, parent Command) Command {
	return cf.FnNew(cmdline, parent)
}

// Factories for each type of command

var CallCommandF CommandFactory[CallOpts] = CommandFactory[CallOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewCallCommand(cmdline, parent) }}
var ConfigCommandF CommandFactory[ConfigOpts] = CommandFactory[ConfigOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewConfigCommand(cmdline, parent) } }
var ConfigGetCommandF CommandFactory[ConfigGetOpts] = CommandFactory[ConfigGetOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewConfigGetCommand(cmdline, parent) } }
var ConfigSetCommandF CommandFactory[ConfigSetOpts] = CommandFactory[ConfigSetOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewConfigSetCommand(cmdline, parent) } }
var ConfigShowCommandF CommandFactory[ConfigShowOpts] = CommandFactory[ConfigShowOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewConfigShowCommand(cmdline, parent) } }
var CreateTwoNodeClusterCommandF CommandFactory[CreateTwoNodeClusterOpts] = CommandFactory[CreateTwoNodeClusterOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewCreateTwoNodeClusterCommand(cmdline, parent) } }
var GenEtcdRunScriptCommandF CommandFactory[GenEtcdRunScriptOpts] = CommandFactory[GenEtcdRunScriptOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewGenEtcdRunScriptCommand(cmdline, parent) } }
var GetCommandF CommandFactory[GetOpts] = CommandFactory[GetOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewGetCommand(cmdline, parent) }}
var HostCommandF CommandFactory[HostOpts] = CommandFactory[HostOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewHostCommand(cmdline, parent) } }
var HostInitCommandF CommandFactory[HostInitOpts] = CommandFactory[HostInitOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewHostInitCommand(cmdline, parent) } }
var InitCommandF CommandFactory[InitOpts] = CommandFactory[InitOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewInitCommand(cmdline, parent) }}
var ListCommandF CommandFactory[ListOpts] = CommandFactory[ListOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewListCommand(cmdline, parent) }}
var LookupCommandF CommandFactory[LookupOpts] = CommandFactory[LookupOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewLookupCommand(cmdline, parent) } }
var MainCommandF CommandFactory[MainOpts] = CommandFactory[MainOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewMainCommand(cmdline, parent) } }
var MiscCommandF CommandFactory[MiscOpts] = CommandFactory[MiscOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewMiscCommand(cmdline, parent) } }
var StatusCommandF CommandFactory[StatusOpts] = CommandFactory[StatusOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewStatusCommand(cmdline, parent) } }
var VmAddCommandF CommandFactory[VmAddOpts] = CommandFactory[VmAddOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewVmAddCommand(cmdline, parent) } }
var VmAllCommandF CommandFactory[VmAllOpts] = CommandFactory[VmAllOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewVmAllCommand(cmdline, parent) } }
var VmCommandF CommandFactory[VmOpts] = CommandFactory[VmOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewVmCommand(cmdline, parent) } }
var VmDelCommandF CommandFactory[VmDelOpts] = CommandFactory[VmDelOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewVmDelCommand(cmdline, parent) } }
var VmGetCommandF CommandFactory[VmGetOpts] = CommandFactory[VmGetOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewVmGetCommand(cmdline, parent) } }
var VmListCommandF CommandFactory[VmListOpts] = CommandFactory[VmListOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewVmListCommand(cmdline, parent) } }
var VmSetCommandF CommandFactory[VmSetOpts] = CommandFactory[VmSetOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewVmSetCommand(cmdline, parent) } }
var WatchCommandF CommandFactory[WatchOpts] = CommandFactory[WatchOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewWatchCommand(cmdline, parent) } }
var WatchScriptCommandF CommandFactory[WatchScriptOpts] = CommandFactory[WatchScriptOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewWatchScriptCommand(cmdline, parent) } }
var WatchSvcCommandF CommandFactory[WatchSvcOpts] = CommandFactory[WatchSvcOpts]{ FnNew: func (cmdline Cmdline, parent Command) Command { return NewWatchSvcCommand(cmdline, parent) } }
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
