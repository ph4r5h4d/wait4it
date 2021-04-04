package redis

import (
	"wait4it/pkg/model"
)

type RedisConnection struct {
	Host          string
	Port          int
	Password      string
	Database      int
	OperationMode string
}

// NewChecker creates a new checker
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
	check := &RedisConnection{}
	check.BuildContext(*c)
	if err := check.Validate(); err != nil {
		return nil, err
	}

	return check, nil
}
