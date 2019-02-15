package api

import (
	"context"
)

type SubmitService struct {
	client *Client
}

func (submit SubmitService) Submit(ctx context.Context, problemId, language, sourceCode string) error {
	var body = struct {
		ProblemId  string `json:"problemId"`
		Language   string `json:"language"`
		SourceCode string `json:"sourceCode"`
	}{problemId, language, sourceCode}

	req, err := submit.client.newRequest(ctx, apiEndpoint, "POST", "submissions", body)
	if err != nil {
		return err
	}

	if err := submit.client.do(req, nil); err != nil {
		return err
	}

	return nil
}
