package commands

import (
	"context"
	"errors"

	"github.com/ken-tunc/aojtool/util"

	"github.com/ken-tunc/aojtool/api"

	"github.com/spf13/cobra"
)

var SubmitLanguage string

var submitCmd = &cobra.Command{
	Use:   "submit [problem-id] [source-file]",
	Short: "Submit a source code.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires at least two args")
		}
		formalLang, err := util.FormalLanguage(SubmitLanguage)
		if SubmitLanguage != "" && err != nil {
			return err
		} else {
			SubmitLanguage = formalLang
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
		loggedIn, err := client.Auth.IsLoggedIn(ctx)
		if err != nil {
			abort(err)
		}

		if !loggedIn {
			cmd.Println("You need to login.")
			return
		}

		user, err := loadUser()
		if err != nil {
			abort(err)
		}

		if SubmitLanguage == "" {
			SubmitLanguage = user.DefaultProgrammingLanguage
		}

		sourceCode, err := util.ReadFile(sourceFile)
		if err != nil {
			abort(err)
		}

		ctx = context.Background()
		if err = client.Submit.Submit(ctx, problemId, SubmitLanguage, sourceCode); err != nil {
			abort(err)
		}
	},
}

func init() {
	submitCmd.Flags().StringVarP(&SubmitLanguage, "language", "l", "", "programming language written in")
	rootCmd.AddCommand(submitCmd)
}
