package mysql

import (
	"wait4it/pkg/model"
)

type MySQLConnection struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
}

// NewChecker creates a new checker
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
	check := &MySQLConnection{}
	check.BuildContext(*c)
	if err := check.Validate(); err != nil {
		return nil, err
	}

	return check, nil
}
