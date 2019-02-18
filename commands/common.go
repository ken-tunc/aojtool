package commands

import (
	"os"
	"path/filepath"

	"github.com/ken-tunc/aojtool/models"
	"github.com/ken-tunc/aojtool/util"
)

var userCache = filepath.Join(util.CacheDir, "user")

func saveUser(user models.User) error {
	return util.SaveData(userCache, &user)
}

func loadUser() (*models.User, error) {
	var user models.User

	err := util.LoadData(userCache, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func removeUser() error {
	return util.RemoveData(userCache)
}

func abort(err error) {
	rootCmd.Println(err)
	os.Exit(1)
}
