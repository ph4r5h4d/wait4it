package tcp

import "testing"

func TestValidate_ValidPort(t *testing.T) {
	c := &Check{
		Addr: "localhost",
		Port: 8080,
	}
	if err := c.Validate(); err != nil {
		t.Errorf("Validate() returned unexpected error: %v", err)
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
			c := &Check{
				Addr: "localhost",
				Port: tt.port,
			}
			if err := c.Validate(); err == nil {
				t.Errorf("Validate() should return error for port %d", tt.port)
			}
		})
	}
}

func TestValidate_BoundaryPorts(t *testing.T) {
	tests := []struct {
		name string
		port int
	}{
		{"min port 1", 1},
		{"max port 65535", 65535},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Check{
				Addr: "localhost",
				Port: tt.port,
			}
			if err := c.Validate(); err != nil {
				t.Errorf("Validate() returned unexpected error for port %d: %v", tt.port, err)
			}
		})
	}
}