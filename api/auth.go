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
	}{id, password}

	req, err := auth.client.newRequest(ctx, apiEndpoint, "POST", "/session", body)
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

	resp, err := auth.client.do(req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("logout failed, status code: %s", resp.Status)
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
