package commands

import (
	"os"
	"path/filepath"

	"github.com/ken-tunc/aojtool/models"
	"github.com/ken-tunc/aojtool/util"
)

var userCache = filepath.Join(util.CacheDir, "user")

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

	data, err := util.ReadBytes(userCache)
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

func abort(err error) {
	rootCmd.Println(err)
	os.Exit(1)
}
