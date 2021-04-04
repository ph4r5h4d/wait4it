package http

import (
	"wait4it/pkg/model"
)

type HttpCheck struct {
	Url    string
	Status int
	Text   string
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
