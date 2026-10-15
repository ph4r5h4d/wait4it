package nats

import "testing"

func TestValidate_Valid(t *testing.T) {
	c := &NatsChecker{
		Host: "localhost",
		Port: 4222,
	}
	if err := c.Validate(); err != nil {
		t.Errorf("Validate() returned unexpected error: %v", err)
	}
}

func TestValidate_EmptyHost(t *testing.T) {
	c := &NatsChecker{
		Host: "",
		Port: 4222,
	}
	if err := c.Validate(); err == nil {
		t.Error("Validate() should return error for empty host")
	}
}

func TestValidate_InvalidPort(t *testing.T) {
	tests := []struct {
		name string
		port int
	}{
		{"zero", 0},
		{"negative", -1},
		{"too high", 65536},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &NatsChecker{
				Host: "localhost",
				Port: tt.port,
			}
			if err := c.Validate(); err == nil {
				t.Errorf("Validate() should return error for port %d", tt.port)
			}
		})
	}
}

func TestBuildURL(t *testing.T) {
	tests := []struct {
		name string
		c    NatsChecker
		want string
	}{
		{
			name: "plain nats",
			c:    NatsChecker{Host: "127.0.0.1", Port: 4222},
			want: "nats://127.0.0.1:4222",
		},
		{
			name: "tls",
			c:    NatsChecker{Host: "nats.example.com", Port: 4443, TLS: true},
			want: "tls://nats.example.com:4443",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.buildURL(); got != tt.want {
				t.Errorf("buildURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTLSEnabled(t *testing.T) {
	tests := []struct {
		ssl  string
		want bool
	}{
		{"disable", false},
		{"", false},
		{"enable", true},
		{"require", true},
		{"true", true},
		{"1", true},
	}

	for _, tt := range tests {
		t.Run(tt.ssl, func(t *testing.T) {
			if got := tlsEnabled(tt.ssl); got != tt.want {
				t.Errorf("tlsEnabled(%q) = %v, want %v", tt.ssl, got, tt.want)
			}
		})
	}
}
