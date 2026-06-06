package tcp

import "testing"

func TestIsPortInValidRange_ValidPorts(t *testing.T) {
	tests := []struct {
		name string
		port int
	}{
		{"min port 1", 1},
		{"common port 80", 80},
		{"common port 443", 443},
		{"common port 3306", 3306},
		{"common port 6379", 6379},
		{"max port 65535", 65535},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Check{Port: tt.port}
			if !c.IsPortInValidRange() {
				t.Errorf("IsPortInValidRange(%d) = false, want true", tt.port)
			}
		})
	}
}

func TestIsPortInValidRange_InvalidPorts(t *testing.T) {
	tests := []struct {
		name string
		port int
	}{
		{"zero", 0},
		{"negative", -1},
		{"too high 65536", 65536},
		{"very negative", -1000},
		{"large number", 100000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Check{Port: tt.port}
			if c.IsPortInValidRange() {
				t.Errorf("IsPortInValidRange(%d) = true, want false", tt.port)
			}
		})
	}
}