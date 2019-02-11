package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ken-tunc/aojtool/util"

	"github.com/ken-tunc/aojtool/models"
)

var userCache = filepath.Join(util.CacheDir, "user")

type AuthService struct {
	client *Client
}

func (auth AuthService) Login(ctx context.Context, id, password string) (*models.User, error) {
	var body = struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}{ID: id, Password: password}

	req, err := auth.client.newRequest(ctx, "POST", "/session", body)
	if err != nil {
		return nil, err
	}

	var user models.User
	resp, err := auth.client.do(req, &user)
	if err != nil {
		return nil, fmt.Errorf("login failed")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("login failed, status code: %s", resp.Status)
	}

	err = auth.client.SaveCookies()
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (auth AuthService) Logout(ctx context.Context) error {
	req, err := auth.client.newRequest(ctx, "DELETE", "/session", nil)
	if err != nil {
		return err
	}

	resp, err := auth.client.do(req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("logout failed, status code: %s", resp.Status)
	}

	err = auth.client.RemoveCookies()
	if err != nil {
		return err
	}

	return nil
}

func (auth AuthService) IsLoggedIn(ctx context.Context) (bool, error) {
	req, err := auth.client.newRequest(ctx, "GET", "/self", nil)
	if err != nil {
		return false, err
	}

	resp, err := auth.client.do(req, nil)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == 200 {
		return true, nil
	} else {
		return false, nil
	}
}

func (auth AuthService) SaveUser(user models.User) error {
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

func (auth AuthService) MaybeLoadUser() (*models.User, error) {
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

func (auth AuthService) RemoveUser() error {
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
