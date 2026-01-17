package dns

import "errors"

var (
	ErrEmptyHost         = errors.New("hostname is required for DNS check")
	ErrInvalidRecordType = errors.New("invalid DNS record type, must be one of: A, AAAA, CNAME, MX, TXT, SRV, NS, PTR")
	ErrNoRecords         = errors.New("no DNS records found")
	ErrExpectedNotFound  = errors.New("expected value not found in DNS records")
)
