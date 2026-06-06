package mongodb

import "testing"

func TestValidate_ValidNoAuth(t *testing.T) {
	m := &MongoDbConnection{
		Host: "localhost",
		Port: 27017,
	}
	if err := m.Validate(); err != nil {
		t.Errorf("Validate() returned unexpected error: %v", err)
	}
}

func TestValidate_ValidWithAuth(t *testing.T) {
	m := &MongoDbConnection{
		Host:     "localhost",
		Port:     27017,
		Username: "admin",
		Password: "secret",
	}
	if err := m.Validate(); err != nil {
		t.Errorf("Validate() returned unexpected error: %v", err)
	}
}

func TestValidate_EmptyHost(t *testing.T) {
	m := &MongoDbConnection{
		Host: "",
		Port: 27017,
	}
	if err := m.Validate(); err == nil {
		t.Error("Validate() should return error for empty host")
	}
}

func TestValidate_UsernameWithoutPassword(t *testing.T) {
	m := &MongoDbConnection{
		Host:     "localhost",
		Port:     27017,
		Username: "admin",
		Password: "",
	}
	if err := m.Validate(); err == nil {
		t.Error("Validate() should return error when username is set but password is empty")
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
			m := &MongoDbConnection{
				Host: "localhost",
				Port: tt.port,
			}
			if err := m.Validate(); err == nil {
				t.Errorf("Validate() should return error for port %d", tt.port)
			}
		})
	}
}