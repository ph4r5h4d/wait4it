package mysql

import "testing"

func TestBuildConnectionString_Full(t *testing.T) {
	m := MySQLConnection{
		Host:         "localhost",
		Port:         3306,
		Username:     "root",
		Password:     "secret",
		DatabaseName: "testdb",
	}
	got := m.BuildConnectionString()
	want := "root:secret@tcp(localhost:3306)/testdb"
	if got != want {
		t.Errorf("BuildConnectionString() = %q, want %q", got, want)
	}
}

func TestBuildConnectionString_NoPassword(t *testing.T) {
	m := MySQLConnection{
		Host:         "localhost",
		Port:         3306,
		Username:     "root",
		Password:     "",
		DatabaseName: "testdb",
	}
	got := m.BuildConnectionString()
	want := "root@tcp(localhost:3306)/testdb"
	if got != want {
		t.Errorf("BuildConnectionString() = %q, want %q", got, want)
	}
}

func TestBuildConnectionString_NoDatabase(t *testing.T) {
	m := MySQLConnection{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "secret",
	}
	got := m.BuildConnectionString()
	want := "root:secret@tcp(localhost:3306)/"
	if got != want {
		t.Errorf("BuildConnectionString() = %q, want %q", got, want)
	}
}

func TestBuildConnectionString_NoPasswordNoDatabase(t *testing.T) {
	m := MySQLConnection{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
	}
	got := m.BuildConnectionString()
	want := "root@tcp(localhost:3306)/"
	if got != want {
		t.Errorf("BuildConnectionString() = %q, want %q", got, want)
	}
}

func TestBuildConnectionString_CustomHostAndPort(t *testing.T) {
	m := MySQLConnection{
		Host:     "db.example.com",
		Port:     3307,
		Username: "admin",
		Password: "p@ss:w0rd",
	}
	got := m.BuildConnectionString()
	want := "admin:p@ss:w0rd@tcp(db.example.com:3307)/"
	if got != want {
		t.Errorf("BuildConnectionString() = %q, want %q", got, want)
	}
}
