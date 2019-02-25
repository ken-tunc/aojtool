package commands

import (
	"context"
	"fmt"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

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

		var id string
		cmd.Print("AOJ user id: ")
		_, err = fmt.Scan(&id)
		if err != nil {
			abort(err)
		}

		cmd.Print("password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			abort(err)
		}
		cmd.Println()
		password := string(bytePassword)

		ctx = context.Background()
		user, err := client.Auth.Login(ctx, id, password)
		if err != nil {
			abort(err)
		}

		if err = saveUser(*user); err != nil {
			abort(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
