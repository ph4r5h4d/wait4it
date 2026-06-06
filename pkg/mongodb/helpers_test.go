package mongodb

import "testing"

func TestBuildConnectionString_WithAuth(t *testing.T) {
	m := &MongoDbConnection{
		Host:     "localhost",
		Port:     27017,
		Username: "admin",
		Password: "secret",
	}
	got := m.buildConnectionString()
	want := "mongodb://admin:secret@localhost:27017"
	if got != want {
		t.Errorf("buildConnectionString() = %q, want %q", got, want)
	}
}

func TestBuildConnectionString_WithoutAuth(t *testing.T) {
	m := &MongoDbConnection{
		Host: "localhost",
		Port: 27017,
	}
	got := m.buildConnectionString()
	want := "mongodb://localhost:27017"
	if got != want {
		t.Errorf("buildConnectionString() = %q, want %q", got, want)
	}
}

func TestBuildConnectionString_EmptyUsername(t *testing.T) {
	m := &MongoDbConnection{
		Host:     "localhost",
		Port:     27017,
		Username: "",
		Password: "secret",
	}
	got := m.buildConnectionString()
	// Empty username means no auth, even if password is set
	want := "mongodb://localhost:27017"
	if got != want {
		t.Errorf("buildConnectionString() = %q, want %q", got, want)
	}
}

func TestBuildConnectionString_CustomHostAndPort(t *testing.T) {
	m := &MongoDbConnection{
		Host:     "db.example.com",
		Port:     27018,
		Username: "user",
		Password: "p@ss",
	}
	got := m.buildConnectionString()
	want := "mongodb://user:p@ss@db.example.com:27018"
	if got != want {
		t.Errorf("buildConnectionString() = %q, want %q", got, want)
	}
}