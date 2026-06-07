//go:build integration

package mongodb

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	tcmongo "github.com/testcontainers/testcontainers-go/modules/mongodb"
	"wait4it/pkg/model"
)

func TestMongoDBIntegration_Check(t *testing.T) {
	ctx := context.Background()

	mongoContainer, err := tcmongo.Run(ctx, "mongo:7")
	if err != nil {
		t.Fatalf("failed to start MongoDB container: %v", err)
	}
	defer func() { _ = mongoContainer.Terminate(ctx) }()

	host, err := mongoContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := mongoContainer.MappedPort(ctx, "27017/tcp")
	if err != nil {
		t.Fatalf("failed to get container port: %v", err)
	}

	portInt, err := strconv.Atoi(port.Port())
	if err != nil {
		t.Fatalf("failed to parse port: %v", err)
	}

	tests := []struct {
		name       string
		host       string
		port       int
		username   string
		password   string
		wantOk     bool
		wantErrMsg string
	}{
		{
			name:     "successful connection without auth",
			host:     host,
			port:     portInt,
			username: "",
			password: "",
			wantOk:   true,
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

			// Retry loop: the Check() method may return (false, false, err) for
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

func TestMongoDBIntegration_WithAuth(t *testing.T) {
	ctx := context.Background()

	// Start MongoDB with auth enabled
	mongoContainer, err := tcmongo.Run(ctx,
		"mongo:7",
		tcmongo.WithUsername("root"),
		tcmongo.WithPassword("rootpassword"),
	)
	if err != nil {
		t.Fatalf("failed to start MongoDB container: %v", err)
	}
	defer func() { _ = mongoContainer.Terminate(ctx) }()

	host, err := mongoContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := mongoContainer.MappedPort(ctx, "27017/tcp")
	if err != nil {
		t.Fatalf("failed to get container port: %v", err)
	}

	portInt, err := strconv.Atoi(port.Port())
	if err != nil {
		t.Fatalf("failed to parse port: %v", err)
	}

	tests := []struct {
		name       string
		host       string
		port       int
		username   string
		password   string
		wantOk     bool
		wantErrMsg string
	}{
		{
			name:     "successful connection with auth",
			host:     host,
			port:     portInt,
			username: "root",
			password: "rootpassword",
			wantOk:   true,
		},
		{
			name:     "wrong password",
			host:     host,
			port:     portInt,
			username: "root",
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

			// Retry loop: the Check() method may return (false, false, err) for
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

func TestMongoDBIntegration_Validate(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		port     int
		username string
		password string
		wantErr  bool
	}{
		{
			name:    "valid config without auth",
			host:    "localhost",
			port:    27017,
			wantErr: false,
		},
		{
			name:    "valid config with auth",
			host:    "localhost",
			port:    27017,
			username: "root",
			password: "rootpassword",
			wantErr: false,
		},
		{
			name:    "empty host",
			host:    "",
			port:    27017,
			wantErr: true,
		},
		{
			name:     "username without password",
			host:     "localhost",
			port:     27017,
			username: "root",
			password: "",
			wantErr:  true,
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
				Username:       tt.username,
				PasswordValue:  tt.password,
			}

			_, err := NewChecker(cx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChecker() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}