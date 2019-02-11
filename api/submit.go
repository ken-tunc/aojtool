package api

import (
	"context"
	"fmt"
)

type SubmitService struct {
	client *Client
}

func (submit SubmitService) Submit(ctx context.Context, problemId, language, sourceCode string) error {
	if err := submit.client.setEndpoint(apiEndpoint); err != nil {
		return err
	}

	var body = struct {
		ProblemId  string `json:"problemId"`
		Language   string `json:"language"`
		SourceCode string `json:"sourceCode"`
	}{problemId, language, sourceCode}

	req, err := submit.client.newRequest(ctx, "POST", "submissions", body)
	if err != nil {
		return err
	}

	resp, err := submit.client.do(req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("login failed, status code: %s", resp.Status)
	}

	return nil
}
