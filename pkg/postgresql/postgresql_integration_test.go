//go:build integration

package postgresql

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"wait4it/pkg/model"
)

func TestPostgreSQLIntegration_Check(t *testing.T) {
	ctx := context.Background()

	postgresContainer, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("postgres"),
	)
	if err != nil {
		t.Fatalf("failed to start PostgreSQL container: %v", err)
	}
	defer func() { _ = postgresContainer.Terminate(ctx) }()

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432/tcp")
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
		database   string
		sslmode    string
		wantOk     bool
		wantErrMsg string
	}{
		{
			name:     "successful connection",
			host:     host,
			port:     portInt,
			username: "postgres",
			password: "postgres",
			database: "testdb",
			sslmode:  "disable",
			wantOk:   true,
		},
		{
			name:       "wrong password",
			host:       host,
			port:       portInt,
			username:   "postgres",
			password:   "wrongpassword",
			database:   "testdb",
			sslmode:    "disable",
			wantOk:     false,
			wantErrMsg: "",
		},
		{
			name:       "wrong database",
			host:       host,
			port:       portInt,
			username:   "postgres",
			password:   "postgres",
			database:   "nonexistent",
			sslmode:    "disable",
			wantOk:     false,
			wantErrMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cx := &model.CheckContext{
				Host:         tt.host,
				Port:         tt.port,
				Username:      tt.username,
				PasswordValue: tt.password,
				DatabaseName:  tt.database,
				DBConf: model.DatabaseSpecificConf{
					SSLMode: tt.sslmode,
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

func TestPostgreSQLIntegration_Validate(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		port    int
		username string
		wantErr bool
	}{
		{
			name:    "valid config",
			host:    "localhost",
			port:    5432,
			username: "postgres",
			wantErr: false,
		},
		{
			name:    "empty host",
			host:    "",
			port:    5432,
			username: "postgres",
			wantErr: true,
		},
		{
			name:    "empty username",
			host:    "localhost",
			port:    5432,
			username: "",
			wantErr: true,
		},
		{
			name:    "invalid port",
			host:    "localhost",
			port:    99999,
			username: "postgres",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cx := &model.CheckContext{
				Host:         tt.host,
				Port:         tt.port,
				Username:      tt.username,
				PasswordValue: "postgres",
				DatabaseName:  "testdb",
				DBConf: model.DatabaseSpecificConf{
					SSLMode: "disable",
				},
			}

			_, err := NewChecker(cx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChecker() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}