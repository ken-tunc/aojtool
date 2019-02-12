package commands

import (
	"context"
	"fmt"
	"text/tabwriter"

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
		cmd.Println("[AOJ user status]")
		printUser(cmd, user)

		if Size < 1 {
			return
		}

		ctx = context.Background()
		records, err := client.Status.FindSubmissionRecords(ctx, user, Size)
		if err != nil {
			abort(err)
		}

		for i, record := range records {
			cmd.Println()
			if i == 0 {
				cmd.Printf("[Recent %d submission status]\n", Size)
			}
			printSubmissionRecord(cmd, record)
		}
	},
}

func init() {
	statusCmd.Flags().IntVarP(&Size, "size", "n", 5, "the number of displayed submission records")
	rootCmd.AddCommand(statusCmd)
}

func printUser(cmd *cobra.Command, user *models.User) {
	w := tabwriter.NewWriter(cmd.OutOrStderr(), 0, 0, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "User ID\t%s\n", user.ID)
	fmt.Fprintf(w, "Last Submit Date\t%s\n", util.TimeFromUnix(user.LastSubmitDate))
	fmt.Fprintf(w, "Default Programming Language\t%s\n", user.DefaultProgrammingLanguage)
	w.Flush()
}

func printSubmissionRecord(cmd *cobra.Command, record models.SubmissionRecord) {
	w := tabwriter.NewWriter(cmd.OutOrStderr(), 0, 0, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "Judge ID\t%d\n", record.JudgeId)
	fmt.Fprintf(w, "Problem ID\t%s\n", record.ProblemId)
	fmt.Fprintf(w, "Submission Date\t%s\n", util.TimeFromUnix(record.SubmissionDate))
	fmt.Fprintf(w, "Language\t%s\n", record.Language)
	fmt.Fprintf(w, "Status\t%s\n", record.Status.String())
	fmt.Fprintf(w, "Score\t%d\n", record.Score)
	if record.Accuracy != nil {
		fmt.Fprintf(w, "Accuracy\t%s\n", *record.Accuracy)
	}
	w.Flush()
}
