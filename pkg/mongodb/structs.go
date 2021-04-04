package mongodb

import (
	"wait4it/pkg/model"
)

type MongoDbConnection struct {
	Host     string
	Port     int
	Username string
	Password string
}

// NewChecker creates a new checker
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
	check := &MongoDbConnection{}
	check.BuildContext(*c)
	if err := check.Validate(); err != nil {
		return nil, err
	}

	return check, nil
}
