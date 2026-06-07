package dns

import (
	"context"
	"net"
	"wait4it/pkg/model"
)

// Resolver is an interface for DNS lookups, matching the methods we use from net.Resolver.
type Resolver interface {
	LookupIP(ctx context.Context, network, host string) ([]net.IP, error)
	LookupCNAME(ctx context.Context, host string) (string, error)
	LookupMX(ctx context.Context, host string) ([]*net.MX, error)
	LookupTXT(ctx context.Context, host string) ([]string, error)
	LookupSRV(ctx context.Context, service, proto, name string) (string, []*net.SRV, error)
	LookupNS(ctx context.Context, host string) ([]*net.NS, error)
	LookupAddr(ctx context.Context, addr string) ([]string, error)
}

type Check struct {
	Host       string
	RecordType string
	Expected   string
	Server     string
	resolver   Resolver
}

// NewChecker creates a new DNS checker instance
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
	recordType := c.DNSConf.RecordType
	if recordType == "" {
		recordType = "A"
	}

	ch := &Check{
		Host:       c.Host,
		RecordType: recordType,
		Expected:   c.DNSConf.Expected,
		Server:     c.DNSConf.Server,
	}

	if err := ch.Validate(); err != nil {
		return nil, err
	}

	return ch, nil
}

func (d *Check) BuildContext(cx model.CheckContext) {
	d.Host = cx.Host
	d.RecordType = cx.DNSConf.RecordType
	d.Expected = cx.DNSConf.Expected
	d.Server = cx.DNSConf.Server
}

func (d *Check) Validate() error {
	if d.Host == "" {
		return ErrEmptyHost
	}

	switch d.RecordType {
	case "A", "AAAA", "CNAME", "MX", "TXT", "SRV", "NS", "PTR":
		// Valid record types
	default:
		return ErrInvalidRecordType
	}

	return nil
}
