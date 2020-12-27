package memcached

import (
	"wait4it/pkg/model"
)

type MemcachedConnection struct {
	Host string
	Port int
}

// NewChecker creates a new checker
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
	check := &MemcachedConnection{}
	check.BuildContext(*c)
	if err := check.Validate(); err != nil {
		return nil, err
	}

	return check, nil
}
