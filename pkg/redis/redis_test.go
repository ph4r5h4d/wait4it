package redis

import "testing"

func TestValidate_ValidStandalone(t *testing.T) {
	m := &RedisConnection{
		Host:          "localhost",
		Port:          6379,
		OperationMode: Standalone,
	}
	if err := m.Validate(); err != nil {
		t.Errorf("Validate() returned unexpected error: %v", err)
	}
}

func TestValidate_ValidCluster(t *testing.T) {
	m := &RedisConnection{
		Host:          "localhost",
		Port:          6379,
		OperationMode: Cluster,
	}
	if err := m.Validate(); err != nil {
		t.Errorf("Validate() returned unexpected error: %v", err)
	}
}

func TestValidate_EmptyHost(t *testing.T) {
	m := &RedisConnection{
		Host:          "",
		Port:          6379,
		OperationMode: Standalone,
	}
	if err := m.Validate(); err == nil {
		t.Error("Validate() should return error for empty host")
	}
}

func TestValidate_InvalidOperationMode(t *testing.T) {
	m := &RedisConnection{
		Host:          "localhost",
		Port:          6379,
		OperationMode: "invalid",
	}
	if err := m.Validate(); err == nil {
		t.Error("Validate() should return error for invalid operation mode")
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
			m := &RedisConnection{
				Host:          "localhost",
				Port:          tt.port,
				OperationMode: Standalone,
			}
			if err := m.Validate(); err == nil {
				t.Errorf("Validate() should return error for port %d", tt.port)
			}
		})
	}
}