package oracle

import "testing"

func TestValidate_Valid(t *testing.T) {
	o := &OracleConnection{
		Host:     "localhost",
		Port:     1521,
		Username: "system",
	}
	if err := o.Validate(); err != nil {
		t.Errorf("Validate() returned unexpected error: %v", err)
	}
}

func TestValidate_EmptyHost(t *testing.T) {
	o := &OracleConnection{
		Host:     "",
		Port:     1521,
		Username: "system",
	}
	if err := o.Validate(); err == nil {
		t.Error("Validate() should return error for empty host")
	}
}

func TestValidate_EmptyUsername(t *testing.T) {
	o := &OracleConnection{
		Host:     "localhost",
		Port:     1521,
		Username: "",
	}
	if err := o.Validate(); err == nil {
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
			o := &OracleConnection{
				Host:     "localhost",
				Port:     tt.port,
				Username: "system",
			}
			if err := o.Validate(); err == nil {
				t.Errorf("Validate() should return error for port %d", tt.port)
			}
		})
	}
}