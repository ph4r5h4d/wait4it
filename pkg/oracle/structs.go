package oracle

import (
	"wait4it/pkg/model"
)

type OracleConnection struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
}

// NewChecker creates a new checker
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
	check := &OracleConnection{}
	check.BuildContext(*c)
	if err := check.Validate(); err != nil {
		return nil, err
	}

	return check, nil
}
