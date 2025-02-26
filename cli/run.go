package cli

// func Run () int {
// 	rootCmd := createRootCommand()
// 	rootCmd.AddCommand(cmd.CallCommand())
// 	rootCmd.AddCommand(cmd.CreateConfigCommand())
// 	rootCmd.AddCommand(cmd.CreateGetCommand())
// 	rootCmd.AddCommand(cmd.CreateInitCommand())
// 	rootCmd.AddCommand(cmd.CreateListCommand())
// 	rootCmd.AddCommand(cmd.CreateVmCommand())
// 	err := rootCmd.Execute()
// 	if err != nil {
// 		return 1
// 	}
// 	return 0
// }

// func createRootCommand() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use: "dylt",
// 		Short: "dylt core functions",
// 		Long: "CLI for using core daylight (dylt) features",
// 		SilenceUsage: true,
// 	}
// 	return cmd
// }

// func CreatePositionalValidator (i int, validArgs []string) (cobra.PositionalArgs, error) {
// 	fn := func(cmd *cobra.Command, args[] string) error {
// 		arg := args[i]
// 		for _, validArg := range(validArgs) {
// 			if arg == validArg {
// 				return nil
// 			}
// 		}
// 		return fmt.Errorf("%s is not a valid value", arg)
// 	}
// 	return fn, nil
// }
