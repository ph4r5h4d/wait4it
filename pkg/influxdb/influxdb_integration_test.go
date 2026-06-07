//go:build integration

package influxdb

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"wait4it/pkg/model"
)

func TestInfluxDBIntegration_Check(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "influxdb:2-alpine",
		ExposedPorts: []string{"8086/tcp"},
		Env: map[string]string{
			"DOCKER_INFLUXDB_INIT_MODE":      "setup",
			"DOCKER_INFLUXDB_INIT_USERNAME":  "admin",
			"DOCKER_INFLUXDB_INIT_PASSWORD":  "secret123",
			"DOCKER_INFLUXDB_INIT_ORG":       "myorg",
			"DOCKER_INFLUXDB_INIT_BUCKET":    "mybucket",
			"DOCKER_INFLUXDB_INIT_ADMIN_TOKEN": "my-super-secret-token",
		},
		WaitingFor: wait.ForListeningPort("8086/tcp").WithStartupTimeout(120*time.Second),
	}

	influxContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:         true,
	})
	if err != nil {
		t.Fatalf("failed to start InfluxDB container: %v", err)
	}
	defer func() { _ = influxContainer.Terminate(ctx) }()

	host, err := influxContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := influxContainer.MappedPort(ctx, "8086/tcp")
	if err != nil {
		t.Fatalf("failed to get container port: %v", err)
	}

	portInt, err := strconv.Atoi(port.Port())
	if err != nil {
		t.Fatalf("failed to parse port: %v", err)
	}

	tests := []struct {
		name   string
		host   string
		port   int
		token  string
		wantOk bool
	}{
		{
			name:   "successful connection",
			host:   "http://" + host,
			port:   portInt,
			token:  "my-super-secret-token",
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cx := &model.CheckContext{
				Host: tt.host,
				Port: tt.port,
				InfluxConf: model.InfluxConf{
					Token: tt.token,
				},
			}

			checker, err := NewChecker(cx)
			if err != nil {
				t.Fatalf("NewChecker() returned unexpected error: %v", err)
			}

			checkCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			// Retry loop: InfluxDB may need time to be fully ready
			var ok, endOnError bool
			var checkErr error
			deadline := time.Now().Add(30 * time.Second)
			for time.Now().Before(deadline) {
				ok, endOnError, checkErr = checker.Check(checkCtx)
				if ok == tt.wantOk {
					break
				}
				if endOnError {
					break
				}
				time.Sleep(1 * time.Second)
			}

			if ok != tt.wantOk {
				t.Errorf("Check() ok = %v, want %v; err = %v", ok, tt.wantOk, checkErr)
			}
			if tt.wantOk && !endOnError {
				t.Error("Check() endOnError = false, want true for successful connection")
			}
		})
	}
}

func TestInfluxDBIntegration_Validate(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		port    int
		wantErr bool
	}{
		{
			name:    "valid config",
			host:    "http://localhost",
			port:    8086,
			wantErr: false,
		},
		{
			name:    "empty host",
			host:    "",
			port:    8086,
			wantErr: true,
		},
		{
			name:    "invalid port",
			host:    "http://localhost",
			port:    99999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cx := &model.CheckContext{
				Host: tt.host,
				Port: tt.port,
			}

			_, err := NewChecker(cx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChecker() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}