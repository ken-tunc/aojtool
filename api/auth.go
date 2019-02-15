package api

import (
	"context"

	"github.com/ken-tunc/aojtool/util"

	"github.com/ken-tunc/aojtool/models"
)

type AuthService struct {
	client *Client
}

func (auth AuthService) Login(ctx context.Context, id, password string) (*models.User, error) {
	var body = struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}{id, password}

	req, err := auth.client.newRequest(ctx, apiEndpoint, "POST", "/session", body)
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := auth.client.do(req, &user); err != nil {
		return nil, err
	}

	if err = auth.client.SaveCookies(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (auth AuthService) Logout(ctx context.Context) error {
	req, err := auth.client.newRequest(ctx, apiEndpoint, "DELETE", "/session", nil)
	if err != nil {
		return err
	}

	if err := auth.client.do(req, nil); err != nil {
		return err
	}

	if err = auth.client.RemoveCookies(); err != nil {
		return err
	}

	return nil
}

func (auth AuthService) IsLoggedIn(ctx context.Context) (bool, error) {
	req, err := auth.client.newRequest(ctx, apiEndpoint, "GET", "/self", nil)
	if err != nil {
		return false, err
	}

	if err = auth.client.do(req, nil); err != nil {
		_, ok := err.(util.ApiErrors)
		if ok {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}
