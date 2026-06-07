package http

import (
	"net/http"
	"wait4it/pkg/model"
)

// HTTPClient is an interface that matches the methods we use from http.Client.
// This allows injection of test clients in integration tests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HttpCheck struct {
	Url            string
	Status         int
	Text           string
	FollowRedirect bool
	client         HTTPClient
}

// NewChecker creates a new checker
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
	check := &HttpCheck{}
	check.BuildContext(*c)
	if err := check.Validate(); err != nil {
		return nil, err
	}

	return check, nil
}
