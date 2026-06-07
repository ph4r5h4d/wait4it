//go:build integration

package oracle

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

func TestOracleIntegration_Check(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "gvenzl/oracle-free:slim-faststart",
		ExposedPorts: []string{"1521/tcp"},
		Env: map[string]string{
			"ORACLE_PASSWORD":    "secret",
			"APP_USER":           "app",
			"APP_USER_PASSWORD":  "secret",
		},
		WaitingFor: wait.ForListeningPort("1521/tcp").WithStartupTimeout(180*time.Second),
	}

	oracleContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:         true,
	})
	if err != nil {
		t.Fatalf("failed to start Oracle container: %v", err)
	}
	defer func() { _ = oracleContainer.Terminate(ctx) }()

	host, err := oracleContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := oracleContainer.MappedPort(ctx, "1521/tcp")
	if err != nil {
		t.Fatalf("failed to get container port: %v", err)
	}

	portInt, err := strconv.Atoi(port.Port())
	if err != nil {
		t.Fatalf("failed to parse port: %v", err)
	}

	tests := []struct {
		name         string
		host         string
		port         int
		username     string
		password     string
		databaseName string
		wantOk       bool
	}{
		{
			name:         "successful connection",
			host:         host,
			port:         portInt,
			username:     "app",
			password:     "secret",
			databaseName: "FREEPDB1",
			wantOk:       true,
		},
		{
			name:         "wrong password",
			host:         host,
			port:         portInt,
			username:     "app",
			password:     "wrongpassword",
			databaseName: "FREEPDB1",
			wantOk:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cx := &model.CheckContext{
				Host:          tt.host,
				Port:          tt.port,
				Username:       tt.username,
				PasswordValue:  tt.password,
				DatabaseName:   tt.databaseName,
			}

			checker, err := NewChecker(cx)
			if err != nil {
				t.Fatalf("NewChecker() returned unexpected error: %v", err)
			}

			checkCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
			defer cancel()

			// Retry loop: Oracle may need time to be fully ready
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
				time.Sleep(2 * time.Second)
			}

			if ok != tt.wantOk {
				t.Errorf("Check() ok = %v, want %v; err = %v", ok, tt.wantOk, checkErr)
			}
			if tt.wantOk && !endOnError {
				t.Error("Check() endOnError = false, want true for successful connection")
			}
			if !tt.wantOk && checkErr != nil {
				fmt.Printf("expected error for %q: %v\n", tt.name, checkErr)
			}
		})
	}
}

func TestOracleIntegration_Validate(t *testing.T) {
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
			port:     1521,
			username: "app",
			wantErr:  false,
		},
		{
			name:     "empty host",
			host:     "",
			port:     1521,
			username: "app",
			wantErr:  true,
		},
		{
			name:     "empty username",
			host:     "localhost",
			port:     1521,
			username: "",
			wantErr:  true,
		},
		{
			name:     "invalid port",
			host:     "localhost",
			port:     99999,
			username: "app",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cx := &model.CheckContext{
				Host:          tt.host,
				Port:          tt.port,
				Username:       tt.username,
				PasswordValue:  "secret",
				DatabaseName:   "FREEPDB1",
			}

			_, err := NewChecker(cx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChecker() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}