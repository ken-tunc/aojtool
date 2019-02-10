package commands

import (
	"context"

	"github.com/ken-tunc/aojtool/api"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Aizu Online Judge.",
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
			return
		}

		ctx = context.Background()
		err = client.Auth.Logout(ctx)

		if err != nil {
			abort(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
