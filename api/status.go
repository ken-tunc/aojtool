package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ken-tunc/aojtool/models"
)

type StatusService struct {
	Client *Client
}

func (status StatusService) FindSubmissionRecords(ctx context.Context, user *models.User, size int) ([]models.SubmissionRecord, error) {
	path := fmt.Sprintf("/submission_records/users/%s", user.ID)

	req, err := status.Client.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("size", strconv.Itoa(size))
	req.URL.RawQuery = q.Encode()

	records := make([]models.SubmissionRecord, size)
	resp, err := status.Client.do(req, &records)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get submission records, status code: %s", resp.Status)
	}

	return records, nil
}
