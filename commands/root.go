package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aojtool",
	Short: "A CLI tool for Aizu Online Judge.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		abort(err)
	}
}

func abort(err error) {
	rootCmd.Println(err)
	os.Exit(1)
}
