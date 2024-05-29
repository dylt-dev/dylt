The command will end up in `./cli/cmd`. So create a new source file there for the new command. One possible naming convention is to use the same name as the file in `./lib` that implements the feature.

```
$ touch cli/cmd/call.go
```

Use snippets to flesh out the source for the new command.
`cli-new`
(This snippet has some funny behavior. It is trying to generalize importing the `lib` functionality of the cli, so it can be used in the new command. But this import path is not known, making it tricky to correctly populate the snippet fields. It might be best to look at an existing command implementation, and see what they do.

Create an instance of the command and add it to the rootCmd 
```
	rootCmd.AddCommand(cmd.CreateCallCommand())
```

To write the actual implementation, look for a xxx_test.go file in lib. If you find one it should have one more more examples of how to call the command's implementation.