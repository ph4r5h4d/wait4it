package dns

import (
	"testing"

	"wait4it/pkg/model"
)

func TestValidate_ValidRecordTypes(t *testing.T) {
	tests := []struct {
		name       string
		recordType string
	}{
		{"A record", "A"},
		{"AAAA record", "AAAA"},
		{"CNAME record", "CNAME"},
		{"MX record", "MX"},
		{"TXT record", "TXT"},
		{"SRV record", "SRV"},
		{"NS record", "NS"},
		{"PTR record", "PTR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Check{
				Host:       "example.com",
				RecordType: tt.recordType,
			}
			if err := d.Validate(); err != nil {
				t.Errorf("Validate() with record type %q returned unexpected error: %v", tt.recordType, err)
			}
		})
	}
}

func TestValidate_InvalidRecordType(t *testing.T) {
	tests := []struct {
		name       string
		recordType string
	}{
		{"empty string", ""},
		{"lowercase a", "a"},
		{"random string", "INVALID"},
		{"number", "123"},
		{"SOA not supported", "SOA"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Check{
				Host:       "example.com",
				RecordType: tt.recordType,
			}
			if err := d.Validate(); err != ErrInvalidRecordType {
				t.Errorf("Validate() with record type %q: got error %v, want %v", tt.recordType, err, ErrInvalidRecordType)
			}
		})
	}
}

func TestValidate_EmptyHost(t *testing.T) {
	d := &Check{
		Host:       "",
		RecordType: "A",
	}
	if err := d.Validate(); err != ErrEmptyHost {
		t.Errorf("Validate() with empty host: got error %v, want %v", err, ErrEmptyHost)
	}
}

func TestValidate_ValidHost(t *testing.T) {
	d := &Check{
		Host:       "example.com",
		RecordType: "A",
	}
	if err := d.Validate(); err != nil {
		t.Errorf("Validate() returned unexpected error: %v", err)
	}
}

func TestNewChecker_DefaultRecordType(t *testing.T) {
	c := &model.CheckContext{
		Host: "example.com",
		DNSConf: model.DNSConf{
			// RecordType intentionally left empty
		},
	}
	checker, err := NewChecker(c)
	if err != nil {
		t.Fatalf("NewChecker() returned unexpected error: %v", err)
	}
	dnsCheck, ok := checker.(*Check)
	if !ok {
		t.Fatal("NewChecker() did not return a *Check")
	}
	if dnsCheck.RecordType != "A" {
		t.Errorf("Default RecordType = %q, want %q", dnsCheck.RecordType, "A")
	}
}