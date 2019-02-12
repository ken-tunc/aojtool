package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ken-tunc/aojtool/api"

	"github.com/ken-tunc/aojtool/util"

	"github.com/spf13/cobra"
)

var (
	RunLanguage string
	TimeOutSec  int
)

var runCmd = &cobra.Command{
	Use:   "run [-l language] [-t timeout] [problem-id] [source-file]",
	Short: "Run program with sample inputs.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires at least two args")
		}
		if RunLanguage != "" && !util.IsAcceptableLanguage(RunLanguage) {
			return fmt.Errorf("invalid language: %s", RunLanguage)
		}
		if TimeOutSec < 1 {
			return fmt.Errorf("invalid timeout seconds: %d", TimeOutSec)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var problemId = args[0]
		var sourceFile = args[1]

		client, err := api.NewClient()
		if err != nil {
			abort(err)
		}

		ctx := context.Background()
		samples, err := client.Test.FindSamples(ctx, problemId)
		if err != nil {
			abort(err)
		}

		if RunLanguage == "" {
			ctx := context.Background()
			loggedIn, err := client.Auth.IsLoggedIn(ctx)
			if err != nil {
				abort(err)
			}

			if !loggedIn {
				abort(errors.New("not logged in"))
			}

			user, err := maybeLoadUser()
			if err != nil {
				abort(err)
			}
			RunLanguage = user.DefaultProgrammingLanguage
		}

		processor := util.NewProcessor(RunLanguage, sourceFile)
		for _, sample := range samples {
			cmd.Printf("[Sample %d]\n", sample.Serial)
			timeout := time.Duration(TimeOutSec) * time.Second
			out, err := processor.Exec(sample.In, timeout)
			out = strings.TrimSpace(out)
			oracle := strings.TrimSpace(sample.Out)
			if err != nil {
				abort(err)
			}
			if out == oracle {
				cmd.Println("Pass.")
			} else {
				cmd.Println("Wrong answer...")
				cmd.Printf("output: %s\n", out)
				cmd.Printf("Oracle: %s\n", oracle)
				break
			}
		}
	},
}

func init() {
	runCmd.Flags().StringVarP(&RunLanguage, "language", "l", "", "programming language written in")
	runCmd.Flags().IntVarP(&TimeOutSec, "timeout", "t", 5, "execution timeout seconds")
	rootCmd.AddCommand(runCmd)
}
