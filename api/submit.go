package api

import (
	"context"
	"fmt"
)

type SubmitService struct {
	Client *Client
}

func (submit SubmitService) Submit(ctx context.Context, problemId, language, sourceCode string) error {
	var body = struct {
		ProblemId  string `json:"problemId"`
		Language   string `json:"language"`
		SourceCode string `json:"sourceCode"`
	}{problemId, language, sourceCode}

	req, err := submit.Client.newRequest(ctx, "POST", "submissions", body)
	if err != nil {
		return err
	}

	resp, err := submit.Client.do(req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("login failed, status code: %s", resp.Status)
	}

	return nil
}
