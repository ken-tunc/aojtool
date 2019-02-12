package commands

import (
	"errors"
	"fmt"

	"github.com/ken-tunc/aojtool/util"

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
		if !util.IsAcceptableLanguage(RunLanguage) {
			return fmt.Errorf("invalid language: %s", RunLanguage)
		}
		return nil
	},
}

func init() {
	runCmd.Flags().StringVarP(&RunLanguage, "language", "l", "", "programming language written in")
	rootCmd.AddCommand(runCmd)
}
