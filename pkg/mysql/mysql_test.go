package mysql

import "testing"

func TestValidate_Valid(t *testing.T) {
	m := &MySQLConnection{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
	}
	if err := m.Validate(); err != nil {
		t.Errorf("Validate() returned unexpected error: %v", err)
	}
}

func TestValidate_EmptyHost(t *testing.T) {
	m := &MySQLConnection{
		Host:     "",
		Port:     3306,
		Username: "root",
	}
	if err := m.Validate(); err == nil {
		t.Error("Validate() should return error for empty host")
	}
}

func TestValidate_EmptyUsername(t *testing.T) {
	m := &MySQLConnection{
		Host:     "localhost",
		Port:     3306,
		Username: "",
	}
	if err := m.Validate(); err == nil {
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
			m := &MySQLConnection{
				Host:     "localhost",
				Port:     tt.port,
				Username: "root",
			}
			if err := m.Validate(); err == nil {
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
			m := &MySQLConnection{
				Host:     "localhost",
				Port:     tt.port,
				Username: "root",
			}
			if err := m.Validate(); err != nil {
				t.Errorf("Validate() returned unexpected error for port %d: %v", tt.port, err)
			}
		})
	}
}