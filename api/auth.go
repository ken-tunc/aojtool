package api

import (
	"context"
	"fmt"

	"github.com/ken-tunc/aojtool/models"
)

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
		return nil, err
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
