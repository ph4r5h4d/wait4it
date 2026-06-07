//go:build integration

package mysql

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"
	"wait4it/pkg/model"
)

func TestMySQLIntegration_Check(t *testing.T) {
	ctx := context.Background()

	mysqlContainer, err := tcmysql.Run(ctx,
		"mysql:8.0",
		tcmysql.WithDatabase("testdb"),
		tcmysql.WithUsername("root"),
		tcmysql.WithPassword("secret"),
	)
	if err != nil {
		t.Fatalf("failed to start MySQL container: %v", err)
	}
	defer func() { _ = mysqlContainer.Terminate(ctx) }()

	host, err := mysqlContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := mysqlContainer.MappedPort(ctx, "3306/tcp")
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
		wantOk     bool
		wantErrMsg string
	}{
		{
			name:     "successful connection",
			host:     host,
			port:     portInt,
			username: "root",
			password: "secret",
			database: "testdb",
			wantOk:   true,
		},
		{
			name:       "wrong password",
			host:       host,
			port:       portInt,
			username:   "root",
			password:   "wrongpassword",
			database:   "testdb",
			wantOk:     false,
			wantErrMsg: "",
		},
		{
			name:       "wrong database",
			host:       host,
			port:       portInt,
			username:   "root",
			password:   "secret",
			database:   "nonexistent",
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

func TestMySQLIntegration_Validate(t *testing.T) {
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
			port:    3306,
			username: "root",
			wantErr: false,
		},
		{
			name:    "empty host",
			host:    "",
			port:    3306,
			username: "root",
			wantErr: true,
		},
		{
			name:    "empty username",
			host:    "localhost",
			port:    3306,
			username: "",
			wantErr: true,
		},
		{
			name:    "invalid port",
			host:    "localhost",
			port:    99999,
			username: "root",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cx := &model.CheckContext{
				Host:         tt.host,
				Port:         tt.port,
				Username:      tt.username,
				PasswordValue: "secret",
				DatabaseName:  "testdb",
			}

			_, err := NewChecker(cx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChecker() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}