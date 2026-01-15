package dns

import (
	"wait4it/pkg/model"
)

type check struct {
	host       string
	recordType string
	expected   string
	server     string
}

// NewChecker creates a new DNS checker instance
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
	recordType := c.DNSConf.RecordType
	if recordType == "" {
		recordType = "A"
	}

	ch := &check{
		host:       c.Host,
		recordType: recordType,
		expected:   c.DNSConf.Expected,
		server:     c.DNSConf.Server,
	}

	if err := ch.validate(); err != nil {
		return nil, err
	}

	return ch, nil
}

func (d *check) validate() error {
	if d.host == "" {
		return ErrEmptyHost
	}

	switch d.recordType {
	case "A", "AAAA", "CNAME", "MX", "TXT", "SRV", "NS", "PTR":
		// Valid record types
	default:
		return ErrInvalidRecordType
	}

	return nil
}
