//go:build integration

package kafka

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"wait4it/pkg/model"
)

func TestKafkaIntegration_Check(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "apache/kafka:latest",
		ExposedPorts: []string{"9092/tcp"},
		Env: map[string]string{
			"KAFKA_NODE_ID":            "0",
			"KAFKA_PROCESS_ROLES":     "controller,broker",
			"KAFKA_LISTENERS":         "PLAINTEXT://:9092,CONTROLLER://:9093",
			"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP": "CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT",
			"KAFKA_CONTROLLER_QUORUM_VOTERS":       "0@localhost:9093",
			"KAFKA_CONTROLLER_LISTENER_NAMES":      "CONTROLLER",
		},
		WaitingFor: wait.ForListeningPort("9092/tcp").WithStartupTimeout(120*time.Second),
	}

	kafkaContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:         true,
	})
	if err != nil {
		t.Fatalf("failed to start Kafka container: %v", err)
	}
	defer func() { _ = kafkaContainer.Terminate(ctx) }()

	host, err := kafkaContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := kafkaContainer.MappedPort(ctx, "9092/tcp")
	if err != nil {
		t.Fatalf("failed to get container port: %v", err)
	}

	portInt, err := parsePortInt(port.Port())
	if err != nil {
		t.Fatalf("failed to parse port: %v", err)
	}

	tests := []struct {
		name            string
		host            string
		port            int
		connectionType  string
		wantOk          bool
	}{
		{
			name:           "successful connection",
			host:           host,
			port:           portInt,
			connectionType: "tcp",
			wantOk:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cx := &model.CheckContext{
				Host: tt.host,
				Port: tt.port,
				KafkaConf: model.KafkaConf{
					ConnectionType: tt.connectionType,
				},
			}

			checker, err := NewChecker(cx)
			if err != nil {
				t.Fatalf("NewChecker() returned unexpected error: %v", err)
			}

			checkCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			// Retry loop: Kafka may need time to be fully ready
			var ok, endOnError bool
			var checkErr error
			deadline := time.Now().Add(20 * time.Second)
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

func TestKafkaIntegration_Validate(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		port    int
		wantErr bool
	}{
		{
			name:    "valid config",
			host:    "localhost",
			port:    9092,
			wantErr: false,
		},
		{
			name:    "empty host",
			host:    "",
			port:    9092,
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
				KafkaConf: model.KafkaConf{
					ConnectionType: "tcp",
				},
			}

			_, err := NewChecker(cx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChecker() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func parsePortInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}