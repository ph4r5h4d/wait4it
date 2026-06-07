//go:build integration

package redis

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"wait4it/pkg/model"
)

func TestRedisIntegration_Check(t *testing.T) {
	ctx := context.Background()

	redisContainer, err := tcredis.Run(ctx, "redis:7-alpine")
	if err != nil {
		t.Fatalf("failed to start Redis container: %v", err)
	}
	defer func() { _ = redisContainer.Terminate(ctx) }()

	host, err := redisContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := redisContainer.MappedPort(ctx, "6379/tcp")
	if err != nil {
		t.Fatalf("failed to get container port: %v", err)
	}

	portInt, err := strconv.Atoi(port.Port())
	if err != nil {
		t.Fatalf("failed to parse port: %v", err)
	}

	tests := []struct {
		name           string
		host           string
		port           int
		password       string
		operationMode  string
		wantOk         bool
		wantErrMsg     string
	}{
		{
			name:          "successful standalone connection",
			host:          host,
			port:          portInt,
			password:       "",
			operationMode: Standalone,
			wantOk:         true,
		},
		{
			name:          "wrong password",
			host:          host,
			port:          portInt,
			password:       "wrongpassword",
			operationMode: Standalone,
			wantOk:         false,
			wantErrMsg:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cx := &model.CheckContext{
				Host:         tt.host,
				Port:         tt.port,
				PasswordValue: tt.password,
				DBConf: model.DatabaseSpecificConf{
					OperationMode: tt.operationMode,
				},
			}

			checker, err := NewChecker(cx)
			if err != nil {
				t.Fatalf("NewChecker() returned unexpected error: %v", err)
			}

			checkCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			// Retry loop: the Check() method returns (false, false, nil) for
			// transient connection failures, so we retry until success or timeout.
			var ok, endOnError bool
			var checkErr error
			deadline := time.Now().Add(15 * time.Second)
			for time.Now().Before(deadline) {
				ok, endOnError, checkErr = checker.Check(checkCtx)
				if ok == tt.wantOk {
					break
				}
				// If the checker signals a hard error (endOnError=true), don't retry
				if endOnError {
					break
				}
				time.Sleep(500 * time.Millisecond)
			}

			if ok != tt.wantOk {
				t.Errorf("Check() ok = %v, want %v", ok, tt.wantOk)
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

func TestRedisIntegration_Validate(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		port    int
		wantErr bool
	}{
		{
			name:    "valid config",
			host:    "localhost",
			port:    6379,
			wantErr: false,
		},
		{
			name:    "empty host",
			host:    "",
			port:    6379,
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
				Host:          tt.host,
				Port:          tt.port,
				DBConf: model.DatabaseSpecificConf{
					OperationMode: Standalone,
				},
			}

			_, err := NewChecker(cx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChecker() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}