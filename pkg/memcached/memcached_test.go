package memcached

import "testing"

func TestValidate_Valid(t *testing.T) {
	m := &MemcachedConnection{
		Host: "localhost",
		Port: 11211,
	}
	if err := m.Validate(); err != nil {
		t.Errorf("Validate() returned unexpected error: %v", err)
	}
}

func TestValidate_EmptyHost(t *testing.T) {
	m := &MemcachedConnection{
		Host: "",
		Port: 11211,
	}
	if err := m.Validate(); err == nil {
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
			m := &MemcachedConnection{
				Host: "localhost",
				Port: tt.port,
			}
			if err := m.Validate(); err == nil {
				t.Errorf("Validate() should return error for port %d", tt.port)
			}
		})
	}
}