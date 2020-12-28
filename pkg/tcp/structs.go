package tcp

import (
	"wait4it/pkg/model"
)

const (
	minPort = 1
	maxPort = 65535
)

type check struct {
	addr string
	port int
}

// NewChecker creates a new checker
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
	check := &check{
		addr: c.Host,
		port: c.Port,
	}
	if err := check.validate(); err != nil {
		return nil, err
	}

	return check, nil
}
