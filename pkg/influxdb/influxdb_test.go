package influxdb

import "testing"

func TestValidate_Valid(t *testing.T) {
	i := &InfluxDBConnection{
		Host: "localhost",
		Port: 8086,
	}
	if err := i.Validate(); err != nil {
		t.Errorf("Validate() returned unexpected error: %v", err)
	}
}

func TestValidate_EmptyHost(t *testing.T) {
	i := &InfluxDBConnection{
		Host: "",
		Port: 8086,
	}
	if err := i.Validate(); err == nil {
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
			i := &InfluxDBConnection{
				Host: "localhost",
				Port: tt.port,
			}
			if err := i.Validate(); err == nil {
				t.Errorf("Validate() should return error for port %d", tt.port)
			}
		})
	}
}

func TestBuildConnectionURL_Default(t *testing.T) {
	i := InfluxDBConnection{
		Host: "localhost",
		Port: 8086,
	}
	got := i.BuildConnectionURL()
	want := "http://localhost:8086"
	if got != want {
		t.Errorf("BuildConnectionURL() = %q, want %q", got, want)
	}
}

func TestBuildConnectionURL_HttpsScheme(t *testing.T) {
	i := InfluxDBConnection{
		Host: "https://influx.example.com",
		Port: 8086,
	}
	got := i.BuildConnectionURL()
	want := "https://influx.example.com:8086"
	if got != want {
		t.Errorf("BuildConnectionURL() = %q, want %q", got, want)
	}
}

func TestBuildConnectionURL_HttpScheme(t *testing.T) {
	i := InfluxDBConnection{
		Host: "http://influx.example.com",
		Port: 8086,
	}
	got := i.BuildConnectionURL()
	want := "http://influx.example.com:8086"
	if got != want {
		t.Errorf("BuildConnectionURL() = %q, want %q", got, want)
	}
}