package memcached

import "testing"

func TestBuildConnectionString(t *testing.T) {
	m := &MemcachedConnection{
		Host: "localhost",
		Port: 11211,
	}
	got := m.BuildConnectionString()
	want := "localhost:11211"
	if got != want {
		t.Errorf("BuildConnectionString() = %q, want %q", got, want)
	}
}

func TestBuildConnectionString_CustomHostAndPort(t *testing.T) {
	m := &MemcachedConnection{
		Host: "cache.example.com",
		Port: 11212,
	}
	got := m.BuildConnectionString()
	want := "cache.example.com:11212"
	if got != want {
		t.Errorf("BuildConnectionString() = %q, want %q", got, want)
	}
}