package commands

import (
	"errors"

	"github.com/spf13/cobra"
)

var RunLanguage string

var runCmd = &cobra.Command{
	Use:   "run [-l language] [problem-id] [source-file]",
	Short: "Run program with sample inputs.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires at least two args")
		}
		return nil
	},
}

func init() {
	runCmd.Flags().StringVarP(&RunLanguage, "language", "l", "", "programming language written in")
	rootCmd.AddCommand(runCmd)
}
