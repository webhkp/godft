package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(explainCmd)
}

var explainCmd = &cobra.Command{
	Use:   "explain",
	Short: "Explain task flow",
	Long: `Read the configuation files. \
and show the task flow.`,
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
		fmt.Println("explaining data....")

	},
}
