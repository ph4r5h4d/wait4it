//go:build integration

package rabbitmq

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"wait4it/pkg/model"
)

func TestRabbitMQIntegration_Check(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "rabbitmq:3-management-alpine",
		ExposedPorts: []string{"5672/tcp"},
		Env: map[string]string{
			"RABBITMQ_DEFAULT_USER": "guest",
			"RABBITMQ_DEFAULT_PASS": "guest",
		},
		WaitingFor: wait.ForListeningPort("5672/tcp").WithStartupTimeout(120*time.Second),
	}

	rabbitContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:         true,
	})
	if err != nil {
		t.Fatalf("failed to start RabbitMQ container: %v", err)
	}
	defer func() { _ = rabbitContainer.Terminate(ctx) }()

	host, err := rabbitContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := rabbitContainer.MappedPort(ctx, "5672/tcp")
	if err != nil {
		t.Fatalf("failed to get container port: %v", err)
	}

	portInt, err := strconv.Atoi(port.Port())
	if err != nil {
		t.Fatalf("failed to parse port: %v", err)
	}

	tests := []struct {
		name     string
		host     string
		port     int
		username string
		password string
		wantOk   bool
	}{
		{
			name:     "successful connection",
			host:     host,
			port:     portInt,
			username: "guest",
			password: "guest",
			wantOk:   true,
		},
		{
			name:     "wrong password",
			host:     host,
			port:     portInt,
			username: "guest",
			password: "wrongpassword",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cx := &model.CheckContext{
				Host:          tt.host,
				Port:          tt.port,
				Username:       tt.username,
				PasswordValue:  tt.password,
			}

			checker, err := NewChecker(cx)
			if err != nil {
				t.Fatalf("NewChecker() returned unexpected error: %v", err)
			}

			checkCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			// Retry loop: RabbitMQ may need time to be fully ready
			var ok, endOnError bool
			var checkErr error
			deadline := time.Now().Add(15 * time.Second)
			for time.Now().Before(deadline) {
				ok, endOnError, checkErr = checker.Check(checkCtx)
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
			if tt.wantOk && endOnError {
				t.Errorf("Check() endOnError = true, want false for successful connection")
			}
			if !tt.wantOk && checkErr != nil {
				fmt.Printf("expected error for %q: %v\n", tt.name, checkErr)
			}
		})
	}
}

func TestRabbitMQIntegration_Validate(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		port     int
		username string
		wantErr  bool
	}{
		{
			name:     "valid config",
			host:     "localhost",
			port:     5672,
			username: "guest",
			wantErr:  false,
		},
		{
			name:     "empty host",
			host:     "",
			port:     5672,
			username: "guest",
			wantErr:  true,
		},
		{
			name:     "empty username",
			host:     "localhost",
			port:     5672,
			username: "",
			wantErr:  true,
		},
		{
			name:     "invalid port",
			host:     "localhost",
			port:     99999,
			username: "guest",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cx := &model.CheckContext{
				Host:          tt.host,
				Port:          tt.port,
				Username:       tt.username,
				PasswordValue:  "guest",
			}

			_, err := NewChecker(cx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChecker() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}