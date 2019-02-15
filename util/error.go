package util

import (
	"fmt"
	"strings"
)

type ApiError struct {
	ID      int    `json:"id"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (err ApiError) Error() string {
	return fmt.Sprintf("%d %s: %s", err.ID, err.Code, err.Message)
}

type ApiErrors []ApiError

func (apiErrors ApiErrors) Error() string {
	var errStrings []string
	for _, err := range apiErrors {
		errStrings = append(errStrings, err.Error())
	}
	return strings.Join(errStrings, "\n")
}
