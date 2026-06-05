package influxdb

import (
	"wait4it/pkg/model"
)

type InfluxDBConnection struct {
	Host  string
	Port  int
	Token string
	// Org and Bucket are populated at CLI level (see model.CheckContext.InfluxConf)
	// but not used by the current basic /health readiness check. They remain
	// available in CheckContext for future enhancements.
}

// NewChecker creates a new checker
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
	check := &InfluxDBConnection{}
	check.BuildContext(*c)
	if err := check.Validate(); err != nil {
		return nil, err
	}

	return check, nil
}
