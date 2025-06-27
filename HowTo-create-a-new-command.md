## HowTo - Create a New `dylt` CLI Command

The command will end up in `./cli/cmd`. So create a new source file there for the new command. One possible naming convention is to use the same name as the file in `./lib` that implements the feature.

```
$ touch cli/cmd/call.go
$ touch lib/call.go # optional but recommended ... I think
```

Use snippets to flesh out the source for the new command.
`cli-new`
(This snippet has some funny behavior. It is trying to generalize importing the `lib` functionality of the cli, so it can be used in the new command. But this import path is not known, making it tricky to correctly populate the snippet fields. It might be best to look at an existing command implementation, and see what they do.

Create an instance of the command and add it to the rootCmd 
```
	rootCmd.AddCommand(cmd.CreateCallCommand())
```

To write the actual implementation, look for a xxx_test.go file in lib. If you find one it should have one more more examples of how to call the command's implementation.

### Skeleton of command

Most commands follow the same pattern
* type *Name*Command* struct { // fields }
* func (*cmd) New*Name*Command () // Create command object
* func (*cmd) HandleArgs () // Boiler plate to populate command object from cmdline args
* func (*cmd) Run () // Tease args out of cmd objecr & call 'real' Run*** function
* func () Run*Name* () // Package-level 'real' run() function that executes the command
```
type <Name>Command struct {
	// fields
}
```

```
// Create a new command object
func New<Name>Command() *<Name>Command>{
	// create command
	// init flag vars

	return &cmd
}
```

```
// Parse flags validate arg count, and initialize positional parameters
func (cmd *<Name>Command) HandleArgs (args []string) error {
	// parse flags
	// validate arg count (nop - command takes all remaining args, 0 or more)
	// init positional params

	return nil
}
```

```
// Tease all necessary data out of the command object, and call the package-level Run<Name> function
func (cmd *<Name>Command) Run (args []string) error {
	// Tease out necessary data
	// Execute Run<Name> command

	return nil
}

// Execute the command.
// Making this function package-level, with explicit arguments, means it can be tested & called explicitly, not just via `dylt` invocation
func Run<Name> (scriptPath string, scriptArgs[] string) error {
	slog.Debug("RunCall()", "scriptPath", scriptPath, "scriptArgs", scriptArgs)
	// Call lib.RunScript() with script path and args, & output response
	_, s, err := lib.RunScript(scriptPath, scriptArgs)
	if err != nil { return err }
	fmt.Printf("%s\n", s)

	return nil
}

```

```
func NewCallCommand () *CallCommand {
	// create command
	flagSet := flag.NewFlagSet("call", flag.ExitOnError)
	cmd := CallCommand{FlagSet: flagSet}
	// init flag vars
	flagSet.StringVar(&cmd.ScriptPath, "script-path", "/opt/bin/daylight.sh", "script-path")
	
	return &cmd
}

func (cmd *CallCommand) HandleArgs (args []string) error {
	// parse flags
	err := cmd.Parse(args)
	if err != nil { return err }
	// validate arg count (nop - command takes all remaining args, 0 or more)
	cmdArgs := cmd.Args()
	// init positional params
	cmd.ScriptArgs = cmdArgs

	return nil
}

func (cmd *CallCommand) Run (args []string) error {
	slog.Debug("CallCommand.Run()", "args", args)
	// Parse flags & get positional args
	err := cmd.HandleArgs(args)
	if err != nil { return err }
	// Execute command
	err = RunCall(cmd.ScriptPath, cmd.ScriptArgs)
	if err != nil { return err }

	return nil
}

func RunCall (scriptPath string, scriptArgs[] string) error {
	slog.Debug("RunCall()", "scriptPath", scriptPath, "scriptArgs", scriptArgs)
	// Call lib.RunScript() with script path and args, & output response
	_, s, err := lib.RunScript(scriptPath, scriptArgs)
	if err != nil { return err }
	fmt.Printf("%s\n", s)

	return nil
}
```

### Subcommands

Commands that have subcommands typically create an additional package-level function: `create<Name>SubCommand`

Example (from `watch.go`)
```
func createWatchSubCommand(cmdName string) (Command, error) {
	switch cmdName {
	case "script": return NewWatchScriptCommand(), nil
	case "svc": return NewWatchSvcCommand(), nil
	default: return nil, fmt.Errorf("unrecognized command: %s", cmdName)
	}
}
```

Each subcommand object will have its own implementation of the above pattern.

Note that when a command has subcommands, the execution of the command might consist of creating a subcommand and then delegating execution to the new subcommand. Such functions might themselves follow a standard pattern.

```
func RunWatch(subCommand string, subCmdArgs []string) error {
	slog.Debug("RunWatch()", "subCommand", subCommand, "subCmdArgs", subCmdArgs)
	// Create the subcommand and run it
	subCmd, err := createWatchSubCommand(subCommand)
	if err != nil { return err }
	err = subCmd.Run(subCmdArgs)
	if err != nil { return err }

	return nil
}

```

### Adding to main

Creating the command doesn't magically enable the command to be invoked from the CLI. That needs to be added to `main.go#createMainSubCommand()"

```
	switch sCmd {
	// *** Add a line here that follows this idiom for the new command ***
	case "call": return clicmd.NewCallCommand(), nil
	case "config": return clicmd.NewConfigCommand(), nil
	case "get": return clicmd.NewGetCommand(), nil
	case "host": return clicmd.NewHostCommand(), nil
	case "init": return clicmd.NewInitCommand(), nil
	case "list": return clicmd.NewListCommand(), nil
	case "misc": return clicmd.NewMiscCommand(), nil
	case "vm": return clicmd.NewVmCommand(), nil
	case "watch": return clicmd.NewWatchCommand(), nil
	default: {
		var nilPtr *MainCommand = nil
		nilPtr.PrintUsage()
		return nil, fmt.Errorf("unrecognized subcommand: %s", sCmd)
	}
	}
```

### Add to `PrintUsage()`

#### Create string constants

Simplest case - a description of the command, and a short usage combining command name w description
```
// dylt misc
const USG_Misc_Desc = "Miscellaneous collection of commands"
var USG_Misc_Short = createUsageShort("misc", USG_Misc_Desc)
```

More complicated - description, short usage, and full usage @note might be better to just always make these
```
// dylt host init
const USG_Host_Init_Desc = "prepare a host for daylight"
var USG_Host_Init_Short = createUsageShort("host init", USG_Host_Init_Desc)
var USG_Host_Init = createUsage("host init", USG_Host_Init_Desc)
```

Even more complicated - full usage is multiline for subcommands
```
// dylt watch
const USG_Watch_Desc = "watch daylight resource for changes"
var USG_Watch_Short = createUsageShort("watch", USG_Watch_Desc)
var USG_Watch = []string {
	USG_Watch_Script_Short,
	USG_Watch_Svc_Short,
}
```

#### Add short usage to main CLI usage
```
// dylt (main)
var USG_Main = []string {
	// *** add an entry here for the new command ***
	USG_Call_Short,
	USG_Config_Short,
	USG_Get_Short,
	USG_Host_Short,
	USG_Init_Short,
	USG_List_Short,
	USG_Misc_Short,
	USG_Vm_Short,
	USG_Watch_Short,
}
```

#### Add usage to command

@note there's a snippet for this

```
func (cmd *StatusCommand) PrintUsage () {
	fmt.Println()
	fmt.Printf("\t%s\n", USG_Status_Desc)
	fmt.Println()
}
```


### Adding to bash autocompletion

Add the new command to the `cmdsDylt` array

```
cmdsDylt=(
	# Add the command here
    call
    config
    get
    host
    init
    list
    misc
    vm
    watch
)            
```

Add a `case` clause for the new command to `do-dylt()`
```
# dylt
do-dylt () {
	# ...
        case $1 in 
		    # Add an invocation for the new command here
            call) do-dylt-call; return;;
			config) do-dylt-config; return;;
			get) do-dylt-get; return;;
			host) do-dylt-host; return;;
			init) do-dylt-init; return;;
			list) do-dylt-list; return;;
			misc) do-dylt-misc; return;;
			vm) do-dylt-vm; return;;
			watch) do-dylt-watch; return;;
            *) echo MEAT
        esac
    fi
}


```
