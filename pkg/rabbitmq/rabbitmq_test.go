package rabbitmq

import "testing"

func TestValidate_Valid(t *testing.T) {
	rc := &RabbitChecker{
		Host:     "localhost",
		Port:     5672,
		Username: "guest",
	}
	if err := rc.Validate(); err != nil {
		t.Errorf("Validate() returned unexpected error: %v", err)
	}
}

func TestValidate_EmptyHost(t *testing.T) {
	rc := &RabbitChecker{
		Host:     "",
		Port:     5672,
		Username: "guest",
	}
	if err := rc.Validate(); err == nil {
		t.Error("Validate() should return error for empty host")
	}
}

func TestValidate_EmptyUsername(t *testing.T) {
	rc := &RabbitChecker{
		Host:     "localhost",
		Port:     5672,
		Username: "",
	}
	if err := rc.Validate(); err == nil {
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
			rc := &RabbitChecker{
				Host:     "localhost",
				Port:     tt.port,
				Username: "guest",
			}
			if err := rc.Validate(); err == nil {
				t.Errorf("Validate() should return error for port %d", tt.port)
			}
		})
	}
}

func TestBuildConnectionString(t *testing.T) {
	rc := &RabbitChecker{
		Host:     "localhost",
		Port:     5672,
		Username: "guest",
		Password: "guest",
	}
	got := rc.buildConnectionString()
	want := "amqp://guest:guest@localhost:5672/"
	if got != want {
		t.Errorf("buildConnectionString() = %q, want %q", got, want)
	}
}