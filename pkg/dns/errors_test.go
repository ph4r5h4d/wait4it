package dns

import "testing"

func TestErrorValues(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		message string
	}{
		{"ErrEmptyHost", ErrEmptyHost, "hostname is required for DNS check"},
		{"ErrInvalidRecordType", ErrInvalidRecordType, "invalid DNS record type, must be one of: A, AAAA, CNAME, MX, TXT, SRV, NS, PTR"},
		{"ErrNoRecords", ErrNoRecords, "no DNS records found"},
		{"ErrExpectedNotFound", ErrExpectedNotFound, "expected value not found in DNS records"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.message {
				t.Errorf("error message = %q, want %q", tt.err.Error(), tt.message)
			}
		})
	}
}