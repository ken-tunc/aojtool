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
	if err := test.client.do(req, &testCases); err != nil {
		return nil, err
	}

	return testCases, nil
}
