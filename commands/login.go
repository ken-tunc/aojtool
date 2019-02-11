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

		userId, password, err := promptIdAndPassword(cmd)
		if err != nil {
			abort(err)
		}

		ctx = context.Background()
		user, err := client.Auth.Login(ctx, userId, password)
		if err != nil {
			abort(err)
		}

		err = client.Auth.SaveUser(*user)
		if err != nil {
			abort(err)
		}
	},
}

func promptIdAndPassword(cmd *cobra.Command) (userId, password string, err error) {
	cmd.Print("AOJ user id: ")
	_, err = fmt.Scan(&userId)
	if err != nil {
		return "", "", err
	}

	cmd.Print("password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	cmd.Println()
	if err != nil {
		return "", "", err
	}

	password = string(bytePassword)
	return
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
