package elasticsearch

import "testing"

func TestBuildConnectionString(t *testing.T) {
	esc := &ElasticSearchChecker{
		Host: "localhost",
		Port: 9200,
	}
	got := esc.BuildConnectionString()
	want := "localhost:9200"
	if got != want {
		t.Errorf("BuildConnectionString() = %q, want %q", got, want)
	}
}

func TestBuildConnectionString_CustomHostAndPort(t *testing.T) {
	esc := &ElasticSearchChecker{
		Host: "es.example.com",
		Port: 9201,
	}
	got := esc.BuildConnectionString()
	want := "es.example.com:9201"
	if got != want {
		t.Errorf("BuildConnectionString() = %q, want %q", got, want)
	}
}