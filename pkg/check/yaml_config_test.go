package check

import (
	"os"
	"path/filepath"
	"testing"
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
