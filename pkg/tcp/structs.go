package tcp

import (
	"wait4it/pkg/model"
)

const (
	minPort = 1
	maxPort = 65535
)

type Check struct {
	Addr string
	Port int
}

// NewChecker creates a new checker
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
	check := &Check{
		Addr: c.Host,
		Port: c.Port,
	}
	if err := check.Validate(); err != nil {
		return nil, err
	}

	return check, nil
}

func (c *Check) BuildContext(cx model.CheckContext) {
	c.Addr = cx.Host
	c.Port = cx.Port
}
