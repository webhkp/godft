package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate task flow configuration",
	Long: `Read the configuration files, \
and validate the configurations.`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Optionally run one of the validators provided by cobra
		// if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
		// 	return err
		// }
		// Run the custom validation logic
		// if myapp.IsValidColor(args[0]) {
		//   return nil
		// }
		// return fmt.Errorf("invalid color specified: %s", args[0])
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("validating configuration....")

	},
}
