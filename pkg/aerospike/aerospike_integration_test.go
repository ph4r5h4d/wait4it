//go:build integration

package aerospike

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"wait4it/pkg/model"
)

func TestAerospikeIntegration_Check(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "aerospike/aerospike-server-enterprise:latest",
		ExposedPorts: []string{"3000/tcp"},
		WaitingFor:   wait.ForListeningPort("3000/tcp").WithStartupTimeout(120*time.Second),
	}

	aeroContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:         true,
	})
	if err != nil {
		t.Fatalf("failed to start Aerospike container: %v", err)
	}
	defer func() { _ = aeroContainer.Terminate(ctx) }()

	host, err := aeroContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := aeroContainer.MappedPort(ctx, "3000/tcp")
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
		wantOk bool
	}{
		{
			name:   "successful connection",
			host:   host,
			port:   portInt,
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cx := &model.CheckContext{
				Host: tt.host,
				Port: tt.port,
			}

			checker, err := NewChecker(cx)
			if err != nil {
				t.Fatalf("NewChecker() returned unexpected error: %v", err)
			}

			// Aerospike Check() doesn't accept context (uses _ context.Context)
			var ok, endOnError bool
			var checkErr error
			deadline := time.Now().Add(15 * time.Second)
			for time.Now().Before(deadline) {
				ok, endOnError, checkErr = checker.Check(context.Background())
				if ok == tt.wantOk {
					break
				}
				if endOnError {
					break
				}
				time.Sleep(500 * time.Millisecond)
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

func TestAerospikeIntegration_Validate(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		port    int
		wantErr bool
	}{
		{
			name:    "valid config",
			host:    "localhost",
			port:    3000,
			wantErr: false,
		},
		{
			name:    "empty host",
			host:    "",
			port:    3000,
			wantErr: true,
		},
		{
			name:    "invalid port",
			host:    "localhost",
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