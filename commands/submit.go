package commands

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/ken-tunc/aojtool/util"

	"github.com/ken-tunc/aojtool/models"

	"github.com/ken-tunc/aojtool/api"

	"github.com/spf13/cobra"
)

var SubmitLanguage string

var submitCmd = &cobra.Command{
	Use:   "submit [-l language] [problem-id] [source-file]",
	Short: "Submit a source code.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires at least two args")
		}
		if SubmitLanguage != "" && !util.IsAcceptableLanguage(SubmitLanguage) {
			return fmt.Errorf("invalid language: %s", SubmitLanguage)
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

		var user *models.User

		if !loggedIn {
			abort(errors.New("not logged in"))
		}

		user, err = maybeLoadUser()
		if err != nil {
			abort(err)
		}

		if SubmitLanguage == "" {
			SubmitLanguage = user.DefaultProgrammingLanguage
		}

		exist, err := util.Exists(sourceFile)
		if err != nil {
			abort(err)
		}

		if !exist {
			abort(fmt.Errorf("source file %s doesn't exist", sourceFile))
		}

		byteSourceCode, err := ioutil.ReadFile(sourceFile)
		if err != nil {
			abort(err)
		}

		sourceCode := string(byteSourceCode)
		ctx = context.Background()
		err = client.Submit.Submit(ctx, problemId, SubmitLanguage, sourceCode)

		if err != nil {
			abort(err)
		}
	},
}

func init() {
	submitCmd.Flags().StringVarP(&SubmitLanguage, "language", "l", "", "programming language written in")
	rootCmd.AddCommand(submitCmd)
}
