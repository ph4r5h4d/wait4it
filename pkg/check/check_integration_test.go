//go:build integration

package check

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

func TestMultiCheckIntegration_RequiredBothSucceed(t *testing.T) {
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

	mysqlHost, err := mysqlContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get mysql host: %v", err)
	}
	mysqlPort, err := mysqlContainer.MappedPort(ctx, "3306/tcp")
	if err != nil {
		t.Fatalf("failed to get mysql port: %v", err)
	}
	mysqlPortInt, err := strconv.Atoi(mysqlPort.Port())
	if err != nil {
		t.Fatalf("failed to parse mysql port: %v", err)
	}

	redisContainer, err := tcredis.Run(ctx, "redis:7-alpine")
	if err != nil {
		t.Fatalf("failed to start Redis container: %v", err)
	}
	defer func() { _ = redisContainer.Terminate(ctx) }()

	redisHost, err := redisContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get redis host: %v", err)
	}
	redisPort, err := redisContainer.MappedPort(ctx, "6379/tcp")
	if err != nil {
		t.Fatalf("failed to get redis port: %v", err)
	}
	redisPortInt, err := strconv.Atoi(redisPort.Port())
	if err != nil {
		t.Fatalf("failed to parse redis port: %v", err)
	}

	yamlContent := fmt.Sprintf(`version: "1"
timeout: 30
checks:
  - name: db
    type: mysql
    host: %s
    port: %d
    username: root
    password: secret
    database: testdb
  - name: cache
    type: redis
    host: %s
    port: %d
`, mysqlHost, mysqlPortInt, redisHost, redisPortInt)

	mcfg, err := LoadConfig("", yamlContent)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	err = RunMultiChecks(ctx, mcfg)
	if err != nil {
		t.Errorf("expected success for two required checks, got error: %v", err)
	}
}

func TestMultiCheckIntegration_RequiredPlusOptionalBad(t *testing.T) {
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

	mysqlHost, err := mysqlContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get mysql host: %v", err)
	}
	mysqlPort, err := mysqlContainer.MappedPort(ctx, "3306/tcp")
	if err != nil {
		t.Fatalf("failed to get mysql port: %v", err)
	}
	mysqlPortInt, err := strconv.Atoi(mysqlPort.Port())
	if err != nil {
		t.Fatalf("failed to parse mysql port: %v", err)
	}

	// Bad redis: wrong port + short timeout per check
	yamlContent := fmt.Sprintf(`version: "1"
timeout: 30
checks:
  - name: db
    type: mysql
    host: %s
    port: %d
    username: root
    password: secret
    database: testdb
    timeout: 30
  - name: cache
    type: redis
    host: 127.0.0.1
    port: 1
    optional: true
    timeout: 2
`, mysqlHost, mysqlPortInt)

	mcfg, err := LoadConfig("", yamlContent)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	err = RunMultiChecks(ctx, mcfg)
	if err != nil {
		t.Errorf("expected overall success (optional failed), got error: %v", err)
	}
}

func TestMultiCheckIntegration_FailFastTrueStops(t *testing.T) {
	ctx := context.Background()

	// Only start redis, make first mysql bad + fail_fast true, should fail without needing second container
	redisContainer, err := tcredis.Run(ctx, "redis:7-alpine")
	if err != nil {
		t.Fatalf("failed to start Redis container: %v", err)
	}
	defer func() { _ = redisContainer.Terminate(ctx) }()

	redisHost, err := redisContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get redis host: %v", err)
	}
	redisPort, err := redisContainer.MappedPort(ctx, "6379/tcp")
	if err != nil {
		t.Fatalf("failed to get redis port: %v", err)
	}
	redisPortInt, err := strconv.Atoi(redisPort.Port())
	if err != nil {
		t.Fatalf("failed to parse redis port: %v", err)
	}

	yamlContent := fmt.Sprintf(`version: "1"
timeout: 5
fail_fast: true
checks:
  - name: badfirst
    type: mysql
    host: 127.0.0.1
    port: 1
    username: root
    password: secret
    database: testdb
    timeout: 2
  - name: cache
    type: redis
    host: %s
    port: %d
`, redisHost, redisPortInt)

	mcfg, err := LoadConfig("", yamlContent)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	err = RunMultiChecks(ctx, mcfg)
	if err == nil {
		t.Error("expected error due to fail_fast on first required failure")
	}
}

func TestMultiCheckIntegration_FailFastFalseContinues(t *testing.T) {
	ctx := context.Background()

	redisContainer, err := tcredis.Run(ctx, "redis:7-alpine")
	if err != nil {
		t.Fatalf("failed to start Redis container: %v", err)
	}
	defer func() { _ = redisContainer.Terminate(ctx) }()

	redisHost, err := redisContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get redis host: %v", err)
	}
	redisPort, err := redisContainer.MappedPort(ctx, "6379/tcp")
	if err != nil {
		t.Fatalf("failed to get redis port: %v", err)
	}
	redisPortInt, err := strconv.Atoi(redisPort.Port())
	if err != nil {
		t.Fatalf("failed to parse redis port: %v", err)
	}

	yamlContent := fmt.Sprintf(`version: "1"
timeout: 5
fail_fast: false
checks:
  - name: badfirst
    type: mysql
    host: 127.0.0.1
    port: 1
    username: root
    password: secret
    database: testdb
    timeout: 2
  - name: cache
    type: redis
    host: %s
    port: %d
    timeout: 5
`, redisHost, redisPortInt)

	mcfg, err := LoadConfig("", yamlContent)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	err = RunMultiChecks(ctx, mcfg)
	if err == nil {
		t.Error("expected error (required failure even with fail_fast false)")
	}
}
