package cmd


// Simple method signature for creating a command
type NewCommandFunc[U Command] func (Cmdline, Command) U

// Simple interface for Command Factory objects, that create a single
// type of Command (eg CallCommand)
//
// U - the type of Command created by this factory
type ICommandFactory[U Command] interface {
	New (Cmdline, Command) U
}

// Simple implementation of ICommandFactory. It has a single field - a 
// reference to a function that returns a type of Command.
// This allows the creation of CommandFactories that are parameterized
// on a Command type and on a factory function for creating that command type.
type CommandFactory[U Command] struct {
	FnNew NewCommandFunc[U]
}

func (cf CommandFactory[U]) New (cmdline Cmdline, parent Command) U {
	return cf.FnNew(cmdline, parent)
}

// Factories for each type of command

var CallCommandF CommandFactory[*CallCommand] = CommandFactory[*CallCommand]{ FnNew: NewCallCommand }
var GetCommandF CommandFactory[*GetCommand] = CommandFactory[*GetCommand]{ FnNew: NewGetCommand }
var InitCommandF CommandFactory[*InitCommand] = CommandFactory[*InitCommand]{ FnNew: NewInitCommand }
var ListCommandF CommandFactory[*ListCommand] = CommandFactory[*ListCommand]{ FnNew: NewListCommand }
var StatusCommandF CommandFactory[*StatusCommand] = CommandFactory[*StatusCommand]{ FnNew: NewStatusCommand }
var ConfigCommandF CommandFactory[*ConfigCommand] = CommandFactory[*ConfigCommand]{ FnNew: NewConfigCommand }
var ConfigGetCommandF CommandFactory[*ConfigGetCommand] = CommandFactory[*ConfigGetCommand]{ FnNew: NewConfigGetCommand }
var ConfigSetCommandF CommandFactory[*ConfigSetCommand] = CommandFactory[*ConfigSetCommand]{ FnNew: NewConfigSetCommand }
var ConfigShowCommandF CommandFactory[*ConfigShowCommand] = CommandFactory[*ConfigShowCommand]{ FnNew: NewConfigShowCommand }
var HostCommandF CommandFactory[*HostCommand] = CommandFactory[*HostCommand]{ FnNew: NewHostCommand }
var HostInitCommandF CommandFactory[*HostInitCommand] = CommandFactory[*HostInitCommand]{ FnNew: NewHostInitCommand }
var MiscCommandF CommandFactory[*MiscCommand] = CommandFactory[*MiscCommand]{ FnNew: NewMiscCommand }
var LookupCommandF CommandFactory[*LookupCommand] = CommandFactory[*LookupCommand]{ FnNew: NewLookupCommand }
var CreateTwoNodeClusterCommandF CommandFactory[*CreateTwoNodeClusterCommand] = CommandFactory[*CreateTwoNodeClusterCommand]{ FnNew: NewCreateTwoNodeClusterCommand }
var GenEtcdRunScriptCommandF CommandFactory[*GenEtcdRunScriptCommand] = CommandFactory[*GenEtcdRunScriptCommand]{ FnNew: NewGenEtcdRunScriptCommand }
var VmCommandF CommandFactory[*VmCommand] = CommandFactory[*VmCommand]{ FnNew: NewVmCommand }
var VmAddCommandF CommandFactory[*VmAddCommand] = CommandFactory[*VmAddCommand]{ FnNew: NewVmAddCommand }
var VmAllCommandF CommandFactory[*VmAllCommand] = CommandFactory[*VmAllCommand]{ FnNew: NewVmAllCommand }
var VmDelCommandF CommandFactory[*VmDelCommand] = CommandFactory[*VmDelCommand]{ FnNew: NewVmDelCommand }
var VmGetCommandF CommandFactory[*VmGetCommand] = CommandFactory[*VmGetCommand]{ FnNew: NewVmGetCommand }
var VmListCommandF CommandFactory[*VmListCommand] = CommandFactory[*VmListCommand]{ FnNew: NewVmListCommand }
var VmSetCommandF CommandFactory[*VmSetCommand] = CommandFactory[*VmSetCommand]{ FnNew: NewVmSetCommand }
var WatchCommandF CommandFactory[*WatchCommand] = CommandFactory[*WatchCommand]{ FnNew: NewWatchCommand }
var WatchScriptCommandF CommandFactory[*WatchScriptCommand] = CommandFactory[*WatchScriptCommand]{ FnNew: NewWatchScriptCommand }
var WatchSvcCommandF CommandFactory[*WatchSvcCommand] = CommandFactory[*WatchSvcCommand]{ FnNew: NewWatchSvcCommand }
var MainCommandF CommandFactory[*MainCommand] = CommandFactory[*MainCommand]{ FnNew: NewMainCommand }
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
