package commands

import (
	"context"

	"github.com/ken-tunc/aojtool/api"
	"github.com/ken-tunc/aojtool/models"
	"github.com/ken-tunc/aojtool/util"
	"github.com/spf13/cobra"
)

var Size int

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print user and recent submission status.",
	Run: func(cmd *cobra.Command, args []string) {
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
			cmd.Println("Not logged in.")
			return
		}

		user, err := maybeLoadUser()
		if err != nil {
			abort(err)
		}
		printUser(cmd, user)

		ctx = context.Background()
		records, err := client.Status.FindSubmissionRecords(ctx, user, Size)
		if err != nil {
			abort(err)
		}

		for _, record := range records {
			cmd.Println()
			printSubmissionRecord(cmd, record)
		}
	},
}

func init() {
	statusCmd.Flags().IntVarP(&Size, "size", "n", 5, "the number of displayed submission records")
	rootCmd.AddCommand(statusCmd)
}

func printUser(cmd *cobra.Command, user *models.User) {
	cmd.Printf("User ID: %s\n", user.ID)
	cmd.Printf("Last Submit Date: %s\n", util.TimeFromUnix(user.LastSubmitDate))
	cmd.Printf("Default Programming Language: %s\n", user.DefaultProgrammingLanguage)
}

func printSubmissionRecord(cmd *cobra.Command, record models.SubmissionRecord) {
	cmd.Printf("Judge ID: %d\n", record.JudgeId)
	cmd.Printf("Problem ID: %s\n", record.ProblemId)
	cmd.Printf("Submission Date: %s\n", util.TimeFromUnix(record.SubmissionDate))
	cmd.Printf("Language: %s\n", record.Language)
	cmd.Printf("Status: %s\n", record.Status.String())
	cmd.Printf("Score: %d\n", record.Score)
	if record.Accuracy != nil {
		cmd.Printf("Accuracy: %s\n", *record.Accuracy)
	}
}
