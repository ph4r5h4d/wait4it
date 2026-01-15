package dns

import (
	"context"
	"net"
	"strings"
	"time"
)

// Check performs the DNS lookup and returns (success, endOnError, error)
func (d *check) Check(ctx context.Context) (bool, bool, error) {
	resolver := d.getResolver()

	var records []string
	var err error

	switch d.recordType {
	case "A":
		records, err = d.lookupA(ctx, resolver)
	case "AAAA":
		records, err = d.lookupAAAA(ctx, resolver)
	case "CNAME":
		records, err = d.lookupCNAME(ctx, resolver)
	case "MX":
		records, err = d.lookupMX(ctx, resolver)
	case "TXT":
		records, err = d.lookupTXT(ctx, resolver)
	case "SRV":
		records, err = d.lookupSRV(ctx, resolver)
	case "NS":
		records, err = d.lookupNS(ctx, resolver)
	case "PTR":
		records, err = d.lookupPTR(ctx, resolver)
	}

	if err != nil {
		// DNS errors are typically temporary, don't end on error
		return false, false, err
	}

	if len(records) == 0 {
		return false, false, ErrNoRecords
	}

	// If no expected value, just check if records exist
	if d.expected == "" {
		return true, false, nil
	}

	// Check if expected value is in any record
	for _, record := range records {
		if strings.Contains(record, d.expected) {
			return true, false, nil
		}
	}

	return false, false, ErrExpectedNotFound
}

// getResolver returns a custom resolver if a DNS server is specified
func (d *check) getResolver() *net.Resolver {
	if d.server == "" {
		return net.DefaultResolver
	}

	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			dialer := net.Dialer{
				Timeout: 5 * time.Second,
			}
			return dialer.DialContext(ctx, "udp", d.server)
		},
	}
}

// lookupA returns IPv4 addresses
func (d *check) lookupA(ctx context.Context, r *net.Resolver) ([]string, error) {
	ips, err := r.LookupIP(ctx, "ip4", d.host)
	if err != nil {
		return nil, err
	}
	var results []string
	for _, ip := range ips {
		results = append(results, ip.String())
	}
	return results, nil
}

// lookupAAAA returns IPv6 addresses
func (d *check) lookupAAAA(ctx context.Context, r *net.Resolver) ([]string, error) {
	ips, err := r.LookupIP(ctx, "ip6", d.host)
	if err != nil {
		return nil, err
	}
	var results []string
	for _, ip := range ips {
		results = append(results, ip.String())
	}
	return results, nil
}

// lookupCNAME returns canonical name
func (d *check) lookupCNAME(ctx context.Context, r *net.Resolver) ([]string, error) {
	cname, err := r.LookupCNAME(ctx, d.host)
	if err != nil {
		return nil, err
	}
	return []string{cname}, nil
}

// lookupMX returns mail exchange records
func (d *check) lookupMX(ctx context.Context, r *net.Resolver) ([]string, error) {
	mxs, err := r.LookupMX(ctx, d.host)
	if err != nil {
		return nil, err
	}
	var results []string
	for _, mx := range mxs {
		results = append(results, mx.Host)
	}
	return results, nil
}

// lookupTXT returns text records
func (d *check) lookupTXT(ctx context.Context, r *net.Resolver) ([]string, error) {
	return r.LookupTXT(ctx, d.host)
}

// lookupSRV returns service records
func (d *check) lookupSRV(ctx context.Context, r *net.Resolver) ([]string, error) {
	// Parse SRV format: _service._proto.name
	// For simplicity, we'll do a direct lookup using the host as-is
	_, addrs, err := r.LookupSRV(ctx, "", "", d.host)
	if err != nil {
		return nil, err
	}
	var results []string
	for _, addr := range addrs {
		results = append(results, addr.Target)
	}
	return results, nil
}

// lookupNS returns nameserver records
func (d *check) lookupNS(ctx context.Context, r *net.Resolver) ([]string, error) {
	nss, err := r.LookupNS(ctx, d.host)
	if err != nil {
		return nil, err
	}
	var results []string
	for _, ns := range nss {
		results = append(results, ns.Host)
	}
	return results, nil
}

// lookupPTR returns reverse DNS records
func (d *check) lookupPTR(ctx context.Context, r *net.Resolver) ([]string, error) {
	return r.LookupAddr(ctx, d.host)
}
