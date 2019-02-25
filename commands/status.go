package commands

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/olekukonko/tablewriter"

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
			cmd.Println("You need to login.")
			return
		}

		user, err := loadUser()
		if err != nil {
			abort(err)
		}
		cmd.Println("[AOJ user status]")
		printUser(cmd.OutOrStderr(), user)

		if Size < 1 {
			return
		}

		ctx = context.Background()
		records, err := client.Status.FindSubmissionRecords(ctx, user, Size)
		if err != nil {
			abort(err)
		}

		cmd.Println()
		printSubmissionRecords(cmd.OutOrStderr(), records)
	},
}

func init() {
	statusCmd.Flags().IntVarP(&Size, "size", "n", 5, "the number of displayed submission records")
	rootCmd.AddCommand(statusCmd)
}

func printUser(w io.Writer, user *models.User) {
	status := []string{
		user.ID,
		util.TimeFromUnix(user.LastSubmitDate).Format("2006-01-02 15:04"),
		user.DefaultProgrammingLanguage,
	}

	table := tablewriter.NewWriter(w)
	table.SetCaption(true, "AOJ user status.")
	table.SetHeader([]string{"User ID", "Last Submit Date", "Default Language"})
	table.Append(status)
	table.Render()
}

func printSubmissionRecords(w io.Writer, records []models.SubmissionRecord) {
	var data [][]string
	for _, record := range records {
		data = append(data, []string{
			strconv.Itoa(record.JudgeId),
			record.ProblemId,
			util.TimeFromUnix(record.SubmissionDate).Format("2006-01-02 15:04"),
			record.Language,
			record.Status.String(),
			strconv.Itoa(record.Score),
		})
	}

	table := tablewriter.NewWriter(w)
	table.SetCaption(true, fmt.Sprintf("Recent %d submission records.", len(records)))
	table.SetHeader([]string{"Judge ID", "Problem ID", "Submission Date", "Language", "Status", "Score"})

	for _, v := range data {
		table.Append(v)
	}

	table.Render()
}
