package dns

import (
	"context"
	"net"
	"strings"
)

// Check performs the DNS lookup and returns (success, endOnError, error)
func (d *Check) Check(ctx context.Context) (bool, bool, error) {
	resolver := d.resolver
	if resolver == nil {
		resolver = d.getResolver()
	}

	var records []string
	var err error

	switch d.RecordType {
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
	if d.Expected == "" {
		return true, false, nil
	}

	// Check if expected value is in any record
	for _, record := range records {
		if strings.Contains(record, d.Expected) {
			return true, false, nil
		}
	}

	return false, false, ErrExpectedNotFound
}

// getResolver returns a custom resolver if a DNS server is specified
func (d *Check) getResolver() *net.Resolver {
	if d.Server == "" {
		return net.DefaultResolver
	}

	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			// DNS resolver is currently hardcoded to use UDP only
			dialer := net.Dialer{}
			return dialer.DialContext(ctx, "udp", d.Server)
		},
	}
}

// lookupA returns IPv4 addresses
func (d *Check) lookupA(ctx context.Context, r Resolver) ([]string, error) {
	ips, err := r.LookupIP(ctx, "ip4", d.Host)
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
func (d *Check) lookupAAAA(ctx context.Context, r Resolver) ([]string, error) {
	ips, err := r.LookupIP(ctx, "ip6", d.Host)
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
func (d *Check) lookupCNAME(ctx context.Context, r Resolver) ([]string, error) {
	cname, err := r.LookupCNAME(ctx, d.Host)
	if err != nil {
		return nil, err
	}
	return []string{cname}, nil
}

// lookupMX returns mail exchange records
func (d *Check) lookupMX(ctx context.Context, r Resolver) ([]string, error) {
	mxs, err := r.LookupMX(ctx, d.Host)
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
func (d *Check) lookupTXT(ctx context.Context, r Resolver) ([]string, error) {
	return r.LookupTXT(ctx, d.Host)
}

// lookupSRV returns service records
func (d *Check) lookupSRV(ctx context.Context, r Resolver) ([]string, error) {
	// Parse SRV format: _service._proto.name
	// For simplicity, we'll do a direct lookup using the host as-is
	_, addrs, err := r.LookupSRV(ctx, "", "", d.Host)
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
func (d *Check) lookupNS(ctx context.Context, r Resolver) ([]string, error) {
	nss, err := r.LookupNS(ctx, d.Host)
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
func (d *Check) lookupPTR(ctx context.Context, r Resolver) ([]string, error) {
	return r.LookupAddr(ctx, d.Host)
}
