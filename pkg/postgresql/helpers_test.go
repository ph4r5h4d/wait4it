package postgresql

import "testing"

func TestBuildConnectionString_Full(t *testing.T) {
	pq := PostgresSQLConnection{
		Host:         "localhost",
		Port:         5432,
		Username:     "postgres",
		Password:     "secret",
		DatabaseName: "testdb",
		SSLMode:      "disable",
	}
	got := pq.BuildConnectionString()
	want := "host=localhost port=5432 user=postgres password=secret sslmode=disable dbname=testdb "
	if got != want {
		t.Errorf("BuildConnectionString() = %q, want %q", got, want)
	}
}

func TestBuildConnectionString_SslModeRequire(t *testing.T) {
	pq := PostgresSQLConnection{
		Host:         "db.example.com",
		Port:         5432,
		Username:     "admin",
		Password:     "p@ss",
		DatabaseName: "mydb",
		SSLMode:      "require",
	}
	got := pq.BuildConnectionString()
	want := "host=db.example.com port=5432 user=admin password=p@ss sslmode=require dbname=mydb "
	if got != want {
		t.Errorf("BuildConnectionString() = %q, want %q", got, want)
	}
}

func TestBuildConnectionString_EmptyFields(t *testing.T) {
	pq := PostgresSQLConnection{
		Host:         "",
		Port:         5432,
		Username:     "",
		Password:     "",
		DatabaseName: "",
		SSLMode:      "",
	}
	got := pq.BuildConnectionString()
	// Should still produce a valid format string even with empty fields
	if got == "" {
		t.Error("BuildConnectionString() should not return empty string")
	}
}