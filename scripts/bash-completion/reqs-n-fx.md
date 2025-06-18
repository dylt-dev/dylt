### Requirements
- subcommands
- flags/options
- flags/options -- long-only
- files
- directories
- enums / collections of literals: static, command
- enums / collections of literals: static, option
- enums / collections of literals : dynamic, command
- enums / collections of literals : dynamic, option
- Option-dependent dynamic sub-options (possibly covered above, but it's the hardest one so it'll be good to have an explicit requirement.

### `dylt` calls for each requirement

####  subcommands

Top level `dylt` should work
```
dylt [TAB][TAB]
```


####  flags/options

Not too many commands have flags. `host init` has uid and gid, so that's pretty good
```
dylt host init - [TAB][TAB]
```


####  flags/options -- long-only
I don't currently have a command that takes long and short options. At least I don't think I do. If I find one or create one, that example should be swapped in.
```
dylt host init  -- [TAB][TAB]
```


####  files
`dylt call` supports the `--script-path` flag to let the caller specify an alternative `daylight.sh` script. This flag can be used to test file completion.

```
dylt call --script-path
```


####  directories
`dylt init` is hardcoded to write to ~/.config/dylt. It would be nice to parameterize this and allow a custom config folder
@note this doesn't exist yet, and will be a slight PITA to support

```
dylt init --config-folder [TAB][TAB]
```


####  enums / collections of literals: static, command

This is no different from supporting subcommands. It'd be nice to have a non-subcommand option but I can't think of any

```
dylt config [TAB][TAB]
```


####  enums / collections of literals: static, option

Right now, commands like `dylt list` have their output format hardcoded. It would be nice to support text and json, or at least to have the flags.
@note functionality doesn't exist

```
dylt list --output [TAB]TAB]
```


####  enums / collections of literals : dynamic, command

`dylt vm get` should dynamically support a list of current virtual machine names

```
dylt vm get [TAB][TAB]
```


####  enums / collections of literals : dynamic, option

@note - no idea

```
```


####  Option-dependent dynamic sub-options (possibly covered above, but it's the hardest one so it'll be good to have an explicit requirement.

@note - no idea

```
```

