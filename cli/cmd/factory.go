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
