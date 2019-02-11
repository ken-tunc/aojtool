package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	"github.com/ken-tunc/aojtool/models"
	"github.com/ken-tunc/aojtool/util"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var userCache = filepath.Join(util.CacheDir, "user")

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

func saveUser(user models.User) error {
	absPath, err := util.EnsurePath(userCache)
	if err != nil {
		return err
	}

	byteUser, err := util.Serialize(&user)
	if err != nil {
		return err
	}

	return util.WriteBytes(byteUser, absPath)
}

func maybeLoadUser() (*models.User, error) {
	exist, err := util.Exists(userCache)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, nil
	}

	var user models.User

	data, err := ioutil.ReadFile(userCache)
	if err != nil {
		return nil, err
	}

	err = util.Deserialize(data, &user)
	return &user, err
}

func removeUser() error {
	exist, err := util.Exists(userCache)
	if err != nil {
		return err
	}

	if exist {
		return os.Remove(userCache)
	} else {
		return nil
	}
}
