package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/webhkp/godft/internal/taskflow"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run [format]",
	Short: "Run all task flows",
	Long: `run process from congiruation. \
Populate data to databases if required.`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Optionally run one of the validators provided by cobra
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}
		
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		taskFlowRunner := taskflow.NewTaskFlowRunner(&args)
		taskFlowRunner.Process()

		fmt.Println("done processing")
	},
}
