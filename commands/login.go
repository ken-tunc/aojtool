package commands

import (
	"context"

	"github.com/ken-tunc/aojtool/api"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Aizu Online Judge.",
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

		if loggedIn {
			return
		}

		userId, password, err := promptIdAndPassword(cmd)
		if err != nil {
			abort(err)
		}

		ctx = context.Background()
		user, err := client.Auth.Login(ctx, userId, password)
		if err != nil {
			abort(err)
		}

		err = saveUser(*user)
		if err != nil {
			abort(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
