package aerospike

import (
	"wait4it/pkg/model"
)

type AerospikeConnection struct {
	Host string
	Port int
}

// NewChecker creates a new checker
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
	check := &AerospikeConnection{}
	check.BuildContext(*c)
	if err := check.Validate(); err != nil {
		return nil, err
	}

	return check, nil
}
