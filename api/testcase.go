package api

import (
	"context"
	"fmt"

	"github.com/ken-tunc/aojtool/models"
)

type TestService struct {
	client *Client
}

func (test *TestService) FindSamples(ctx context.Context, problemId string) ([]models.TestCase, error) {
	path := fmt.Sprintf("/testcases/samples/%s", problemId)
	req, err := test.client.newRequest(ctx, datEndpoint, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	testCases := make([]models.TestCase, 0)
	resp, err := test.client.do(req, &testCases)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get samples, status code: %s", resp.Status)
	}

	return testCases, nil
}
