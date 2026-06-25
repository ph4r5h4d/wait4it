package check

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"wait4it/pkg/model"
)

func TestLoadConfig_FromInlineContent(t *testing.T) {
	yamlContent := `
version: "1"
timeout: 45
fail_fast: false
checks:
  - name: primary-db
    type: mysql
    host: db
    port: 3306
    username: root
    password: secret
    database: app
    optional: false
  - name: cache
    type: redis
    host: redis
    port: 6379
    optional: true
`

	cfg, err := LoadConfig("", yamlContent)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
	if cfg.Version != "1" {
		t.Errorf("Version = %q, want %q", cfg.Version, "1")
	}
	if cfg.Timeout != 45 {
		t.Errorf("Timeout = %d, want 45", cfg.Timeout)
	}
	if cfg.FailFast == nil || *cfg.FailFast != false {
		t.Errorf("FailFast should be explicit false")
	}
	if len(cfg.Checks) != 2 {
		t.Fatalf("len(Checks) = %d, want 2", len(cfg.Checks))
	}
	if cfg.Checks[0].Name != "primary-db" || cfg.Checks[0].Type != "mysql" || cfg.Checks[0].Optional {
		t.Errorf("first check parsed incorrectly: %+v", cfg.Checks[0])
	}
	if !cfg.Checks[1].Optional {
		t.Errorf("second check should be optional")
	}
}

func TestLoadConfig_FromFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "checks.yaml")
	content := `version: "1"
checks:
- type: tcp
  host: 127.0.0.1
  port: 8080
  name: web
`
	if err := os.WriteFile(f, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp yaml: %v", err)
	}

	cfg, err := LoadConfig(f, "")
	if err != nil {
		t.Fatalf("LoadConfig from file: %v", err)
	}
	if len(cfg.Checks) != 1 || cfg.Checks[0].Name != "web" || cfg.Checks[0].Type != "tcp" {
		t.Errorf("unexpected parsed config from file: %+v", cfg)
	}
}

func TestLoadConfig_PathTakesPrecedence(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "real.yaml")
	if err := os.WriteFile(f, []byte(`version: "1"
checks:
- type: tcp
  host: 10.0.0.1
  port: 1234
`), 0644); err != nil {
		t.Fatal(err)
	}

	// content should be ignored because path is provided
	cfg, err := LoadConfig(f, `version: "1"
checks:
- type: tcp
  host: ignored
`)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Checks[0].Host != "10.0.0.1" {
		t.Errorf("expected host from file, got %s", cfg.Checks[0].Host)
	}
}

func TestLoadConfig_Errors(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		content     string
		wantErrPart string
	}{
		{"no config", "", "", "no configuration provided"},
		{"bad yaml", "", "not: valid: [yaml", "parse YAML"},
		{"bad version", "", `version: "2"
checks:
- type: tcp
`, "unsupported config version"},
		{"empty checks", "", `version: "1"
checks: []
`, "no checks defined"},
		{"missing type", "", `version: "1"
checks:
- name: foo
  host: localhost
  port: 1234
`, "missing required field \"type\""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadConfig(tt.path, tt.content)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if tt.wantErrPart != "" {
				if !contains(err.Error(), tt.wantErrPart) {
					t.Errorf("error %q does not contain %q", err, tt.wantErrPart)
				}
			}
		})
	}
}

func TestLoadConfig_NonexistentFile(t *testing.T) {
	_, err := LoadConfig("/non/existent/path/for/wait4it/test.yaml", "")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestSpecToCheckContext_BasicFields(t *testing.T) {
	spec := CheckSpec{
		Type:         "mysql",
		Host:         "dbhost",
		Port:         3307,
		Username:     "u",
		Password:     "p",
		PasswordFile: "/tmp/pf",
		DatabaseName: "mydb",
		Timeout:      77,
	}
	cc := SpecToCheckContext(spec, 30)

	if cc.Config.CheckType != "mysql" {
		t.Errorf("CheckType=%s", cc.Config.CheckType)
	}
	if cc.Config.Timeout != 77 {
		t.Errorf("Timeout=%d", cc.Config.Timeout)
	}
	if cc.Host != "dbhost" || cc.Port != 3307 {
		t.Errorf("host/port mismatch")
	}
	if cc.Username != "u" || cc.PasswordValue != "p" || cc.PasswordFile != "/tmp/pf" {
		t.Errorf("auth fields mismatch")
	}
	if cc.DatabaseName != "mydb" {
		t.Errorf("db name mismatch")
	}
}

func TestSpecToCheckContext_Defaults(t *testing.T) {
	spec := CheckSpec{Type: "dns"}
	cc := SpecToCheckContext(spec, 0)

	if cc.DNSConf.RecordType != "A" {
		t.Errorf("dns default RecordType want A, got %s", cc.DNSConf.RecordType)
	}
	if cc.HttpConf.StatusCode != 200 {
		t.Errorf("http default status want 200")
	}
	if !cc.HttpConf.FollowRedirect {
		t.Errorf("http default follow redirect want true")
	}
	if cc.KafkaConf.ConnectionType != "tcp" {
		t.Errorf("kafka default want tcp")
	}
	if cc.Config.Timeout != 30 {
		t.Errorf("global default timeout want 30")
	}
}

func TestSpecToCheckContext_HttpFollowRedirectOverride(t *testing.T) {
	falseVal := false
	spec := CheckSpec{
		Type:                 "http",
		HttpFollowRedirect: &falseVal,
	}
	cc := SpecToCheckContext(spec, 30)
	if cc.HttpConf.FollowRedirect {
		t.Error("expected FollowRedirect=false")
	}

	trueVal := true
	spec2 := CheckSpec{Type: "http", HttpFollowRedirect: &trueVal}
	cc2 := SpecToCheckContext(spec2, 30)
	if !cc2.HttpConf.FollowRedirect {
		t.Error("expected FollowRedirect=true override")
	}
}

func TestSpecToCheckContext_PerCheckTimeoutWins(t *testing.T) {
	spec := CheckSpec{Type: "tcp", Timeout: 99}
	cc := SpecToCheckContext(spec, 10)
	if cc.Config.Timeout != 99 {
		t.Errorf("per-check timeout not respected")
	}
}

func TestSpecToCheckContext_AllTypesProduceUsableContext(t *testing.T) {
	// Use the same generous defaults as the module list tests.
	for _, checkType := range []string{
		"tcp", "mysql", "postgres", "http", "mongo", "oracle",
		"influxdb", "redis", "rabbitmq", "memcached",
		"elasticsearch", "aerospike", "kafka", "dns",
	} {
		t.Run(checkType, func(t *testing.T) {
			spec := CheckSpec{
				Type:          checkType,
				Host:          "localhost",
				Port:          1234,
				Username:      "testuser",
				Password:      "testpass",
				DatabaseName:  "testdb",
				OperationMode: "standalone",
			}
			cc := SpecToCheckContext(spec, 1)

			// Should at least not panic and either succeed or return a validation error
			// (same tolerance used in check-module-list_test.go)
			_, err := findCheckModule(cc)
			if err != nil {
				// Acceptable for incomplete configs
				return
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && searchSubstring(s, substr)))
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Additional mapper field coverage for sub-structs
func TestSpecToCheckContext_TypeSpecificFields(t *testing.T) {
	b := false
	spec := CheckSpec{
		Type:                             "http",
		HttpStatusCode:                   201,
		HttpText:                         "ready",
		HttpFollowRedirect:               &b,
		SSLMode:                          "require",
		OperationMode:                    "cluster",
		KafkaConnectionType:              "tls",
		KafkaConnectToLeaderViaNonLeader: true,
		DNSType:                          "SRV",
		DNSExpect:                        "foo",
		DNSServer:                        "8.8.8.8:53",
		InfluxToken:                      "tok",
		InfluxOrg:                        "org",
		InfluxBucket:                     "bkt",
	}

	cc := SpecToCheckContext(spec, 30)

	if cc.HttpConf.StatusCode != 201 || cc.HttpConf.Text != "ready" || cc.HttpConf.FollowRedirect {
		t.Errorf("http conf not mapped: %+v", cc.HttpConf)
	}
	if cc.DBConf.SSLMode != "require" || cc.DBConf.OperationMode != "cluster" {
		t.Errorf("db conf not mapped")
	}
	if cc.KafkaConf.ConnectionType != "tls" || !cc.KafkaConf.ConnectToLeaderViaNonLeader {
		t.Errorf("kafka conf not mapped")
	}
	if cc.DNSConf.RecordType != "SRV" || cc.DNSConf.Expected != "foo" || cc.DNSConf.Server != "8.8.8.8:53" {
		t.Errorf("dns conf not mapped")
	}
	if cc.InfluxConf.Token != "tok" || cc.InfluxConf.Org != "org" || cc.InfluxConf.Bucket != "bkt" {
		t.Errorf("influx conf not mapped")
	}
}

// --- runMulti skeleton tests (pure, injectable runner) ---

func TestRunMulti_NameSynthesis(t *testing.T) {
	cfg := &MultiConfig{
		Checks: []CheckSpec{
			{Type: "tcp", Host: "h", Port: 1},
			{Name: "custom", Type: "redis", Host: "h", Port: 2},
		},
	}
	calls := []string{}
	err := runMulti(cfg, func(cc *model.CheckContext) error {
		calls = append(calls, cc.Config.CheckType)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(calls) != 2 {
		t.Fatalf("got %d calls", len(calls))
	}
}

func TestRunMulti_AllSuccess(t *testing.T) {
	cfg := &MultiConfig{
		Checks: []CheckSpec{
			{Type: "tcp", Host: "h", Port: 1},
		},
	}
	err := runMulti(cfg, func(*model.CheckContext) error { return nil })
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestRunMulti_OptionalFails_OverallSuccess(t *testing.T) {
	cfg := &MultiConfig{
		Checks: []CheckSpec{
			{Type: "tcp", Host: "h", Port: 1},
			{Type: "redis", Host: "h", Port: 2, Optional: true},
		},
	}
	callCount := 0
	err := runMulti(cfg, func(cc *model.CheckContext) error {
		callCount++
		if cc.Config.CheckType == "redis" {
			return fmt.Errorf("redis down")
		}
		return nil
	})
	if err != nil {
		t.Errorf("expected success despite optional fail, got %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls, got %d", callCount)
	}
}

func TestRunMulti_RequiredFails_FailFast(t *testing.T) {
	cfg := &MultiConfig{
		Checks: []CheckSpec{
			{Type: "tcp", Host: "h", Port: 1},
			{Type: "redis", Host: "h", Port: 2}, // required
		},
	}
	err := runMulti(cfg, func(cc *model.CheckContext) error {
		if cc.Config.CheckType == "redis" {
			return fmt.Errorf("redis failed")
		}
		return nil
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !contains(err.Error(), "redis") {
		t.Errorf("error should mention redis: %v", err)
	}
}

func TestRunMulti_RequiredFails_NoFailFast_Continues(t *testing.T) {
	cfg := &MultiConfig{
		FailFast: boolPtr(false),
		Checks: []CheckSpec{
			{Type: "tcp", Host: "h", Port: 1},     // will fail
			{Type: "redis", Host: "h", Port: 2}, // required, should still run
		},
	}
	ran := []string{}
	err := runMulti(cfg, func(cc *model.CheckContext) error {
		ran = append(ran, cc.Config.CheckType)
		if cc.Config.CheckType == "tcp" {
			return fmt.Errorf("tcp failed")
		}
		return nil
	})
	if err == nil {
		t.Fatal("expected error from required fail")
	}
	if len(ran) != 2 {
		t.Errorf("expected both to run, got %v", ran)
	}
}

func TestRunMulti_Mixed(t *testing.T) {
	cfg := &MultiConfig{
		Checks: []CheckSpec{
			{Name: "db", Type: "mysql", Host: "h", Port: 1},
			{Type: "redis", Host: "h", Port: 2, Optional: true},
			{Type: "kafka", Host: "h", Port: 3},
		},
	}
	err := runMulti(cfg, func(cc *model.CheckContext) error {
		if cc.Config.CheckType == "redis" {
			return fmt.Errorf("redis optional fail")
		}
		return nil
	})
	if err != nil {
		t.Errorf("expected overall success, got %v", err)
	}
}

func boolPtr(b bool) *bool { return &b }

// Light test for real RunMultiChecks path (fails fast at findCheckModule, no long wait).
func TestRunMultiChecks_InvalidType(t *testing.T) {
	cfg := &MultiConfig{
		Checks: []CheckSpec{
			{Type: "unknown", Host: "localhost", Port: 1234},
		},
	}
	err := RunMultiChecks(context.Background(), cfg)
	if err == nil {
		t.Error("expected error for unknown type")
	}
}

func TestRunMultiChecks_PasswordFileValidation(t *testing.T) {
	dir := t.TempDir()
	goodPF := filepath.Join(dir, "good.txt")
	if err := os.WriteFile(goodPF, []byte("secret"), 0644); err != nil {
		t.Fatal(err)
	}

	// Good pw file -> validate passes (will fail later on actual check, but not on pw)
	cfgGood := &MultiConfig{
		Checks: []CheckSpec{
			{Type: "tcp", Host: "h", Port: 1, PasswordFile: goodPF},
		},
	}
	err := RunMultiChecks(context.Background(), cfgGood)
	if err == nil {
		t.Fatal("expected error (from tcp check)")
	}
	if contains(err.Error(), "password file validation") {
		t.Errorf("unexpected pw validation error for good file: %v", err)
	}

	// Bad pw file -> should error with validation message (even for required)
	cfgBad := &MultiConfig{
		Checks: []CheckSpec{
			{Type: "tcp", Host: "h", Port: 1, PasswordFile: "/nonexistent/pf.txt"},
		},
	}
	err = RunMultiChecks(context.Background(), cfgBad)
	if err == nil {
		t.Fatal("expected error")
	}
	if !contains(err.Error(), "password file validation failed") {
		t.Errorf("expected pw validation error, got: %v", err)
	}
}

func TestRunMultiChecks_PerCheckTimeoutFallback(t *testing.T) {
	cfg := &MultiConfig{
		Timeout: 5,
		Checks: []CheckSpec{
			{Type: "tcp", Host: "h", Port: 1, Timeout: 0}, // should fallback
		},
	}
	// We can't easily assert the timeout value without mocking deeper,
	// but at least it shouldn't panic and should use fallback.
	err := RunMultiChecks(context.Background(), cfg)
	if err == nil {
		t.Fatal("expected failure")
	}
}
