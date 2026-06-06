package redis

import "testing"

func TestBuildConnectionString(t *testing.T) {
	m := RedisConnection{
		Host: "localhost",
		Port: 6379,
	}
	got := m.BuildConnectionString()
	want := "localhost:6379"
	if got != want {
		t.Errorf("BuildConnectionString() = %q, want %q", got, want)
	}
}

func TestBuildConnectionString_CustomHostAndPort(t *testing.T) {
	m := RedisConnection{
		Host: "redis.example.com",
		Port: 6380,
	}
	got := m.BuildConnectionString()
	want := "redis.example.com:6380"
	if got != want {
		t.Errorf("BuildConnectionString() = %q, want %q", got, want)
	}
}