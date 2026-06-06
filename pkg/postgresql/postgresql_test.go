package postgresql

import "testing"

func TestValidate_Valid(t *testing.T) {
	pq := &PostgresSQLConnection{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
	}
	if err := pq.Validate(); err != nil {
		t.Errorf("Validate() returned unexpected error: %v", err)
	}
}

func TestValidate_EmptyHost(t *testing.T) {
	pq := &PostgresSQLConnection{
		Host:     "",
		Port:     5432,
		Username: "postgres",
	}
	if err := pq.Validate(); err == nil {
		t.Error("Validate() should return error for empty host")
	}
}

func TestValidate_EmptyUsername(t *testing.T) {
	pq := &PostgresSQLConnection{
		Host:     "localhost",
		Port:     5432,
		Username: "",
	}
	if err := pq.Validate(); err == nil {
		t.Error("Validate() should return error for empty username")
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
			pq := &PostgresSQLConnection{
				Host:     "localhost",
				Port:     tt.port,
				Username: "postgres",
			}
			if err := pq.Validate(); err == nil {
				t.Errorf("Validate() should return error for port %d", tt.port)
			}
		})
	}
}

func TestValidate_ValidBoundaryPorts(t *testing.T) {
	tests := []struct {
		name string
		port int
	}{
		{"min port 1", 1},
		{"max port 65535", 65535},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := &PostgresSQLConnection{
				Host:     "localhost",
				Port:     tt.port,
				Username: "postgres",
			}
			if err := pq.Validate(); err != nil {
				t.Errorf("Validate() returned unexpected error for port %d: %v", tt.port, err)
			}
		})
	}
}