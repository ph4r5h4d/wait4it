package kafka

import (
	"wait4it/pkg/model"
)

type KafkaConnection struct {
	Host                        string
	Port                        int
	ConnectionType              string
	ConnectToLeaderViaNonLeader bool
}

// NewChecker creates a new checker
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
	check := &KafkaConnection{}
	check.BuildContext(*c)
	if err := check.Validate(); err != nil {
		return nil, err
	}

	return check, nil
}
