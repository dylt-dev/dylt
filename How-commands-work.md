Each `dylt` command has several files associcated it
- Set up the command (flags, usage, subcommands)
- Execute the command itself
- Usage information
- Test to execute the command

Let's look at the status command

`cli/cmd/status.go` - Command object
`lib/status.go` - Command execution
`lib/status_test.go` - Tests that execute the command
`cli/cmd/usage.go` - Usage for the command

It's hard to learn how a command's code works if you have no idea what the 
command does or how to execute it.

Read the command's description (`cli/cmd/status.go`)
Read its usage (`cli/cmd/usage.go`)
Run some tests (`lib/status_test.go`)

Once you have familiarized yourself with the purpose of a command and how to use
the command, it's a good time to start looking at the code for the command -
both the code that constructs the command object, and the code that executes the
command.

#### Constructing the command object

Command object construction follows a common pattern that has a common set of
components

- a type for the command, that embeds `flag.FlagSet`
- a factory function that creates a new command object
- a `HandleArgs()` function that validates the number and value of any arguments
to the command
- a `PrintUsage()` function that prints the usage text from `usage.go` for the 
command
- a `Run` function that takes care boiler plate of parsing the command line,
validating args, and eventually running the command.
- a `Run[Command]()` function that delegates to code in `lib` that executes the command

