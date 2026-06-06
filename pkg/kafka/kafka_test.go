package kafka

import "testing"

func TestValidate_Valid(t *testing.T) {
	c := &KafkaConnection{
		Host: "localhost",
		Port: 9092,
	}
	if err := c.Validate(); err != nil {
		t.Errorf("Validate() returned unexpected error: %v", err)
	}
}

func TestValidate_EmptyHost(t *testing.T) {
	c := &KafkaConnection{
		Host: "",
		Port: 9092,
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
			c := &KafkaConnection{
				Host: "localhost",
				Port: tt.port,
			}
			if err := c.Validate(); err == nil {
				t.Errorf("Validate() should return error for port %d", tt.port)
			}
		})
	}
}