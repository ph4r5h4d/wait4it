package check

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"wait4it/pkg/model"
)

// MultiConfig represents the top-level structure of a wait4it YAML configuration
// file supporting multiple checks.
type MultiConfig struct {
	Version  string      `yaml:"version"`
	Timeout  int         `yaml:"timeout"`
	FailFast *bool       `yaml:"fail_fast"`
	Checks   []CheckSpec `yaml:"checks"`
}

// CheckSpec describes a single check entry inside a multi-check YAML configuration.
type CheckSpec struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Optional bool   `yaml:"optional"`
	Timeout  int    `yaml:"timeout"`

	// Common fields (same semantics as CLI flags / envs)
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	PasswordFile string `yaml:"password_file"`
	DatabaseName string `yaml:"database"`

	// Database specific
	SSLMode       string `yaml:"ssl_mode"`
	OperationMode string `yaml:"operation_mode"`

	// HTTP specific
	HttpStatusCode     int    `yaml:"http_status_code"`
	HttpText           string `yaml:"http_text"`
	HttpFollowRedirect *bool  `yaml:"http_follow_redirect"`

	// Kafka specific
	KafkaConnectionType              string `yaml:"kafka_connection_type"`
	KafkaConnectToLeaderViaNonLeader bool   `yaml:"kafka_connect_to_leader_via_non_leader"`

	// DNS specific
	DNSType   string `yaml:"dns_type"`
	DNSExpect string `yaml:"dns_expect"`
	DNSServer string `yaml:"dns_server"`

	// InfluxDB specific
	InfluxToken  string `yaml:"influx_token"`
	InfluxOrg    string `yaml:"influx_org"`
	InfluxBucket string `yaml:"influx_bucket"`
}

// LoadConfig loads a MultiConfig.
// If path is non-empty, the file at that path is read. Otherwise the provided
// content string is used as the YAML source.
// Path takes precedence when both are supplied.
func LoadConfig(path, content string) (*MultiConfig, error) {
	if path != "" {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
		}
		content = string(b)
	}
	if content == "" {
		return nil, fmt.Errorf("no configuration provided (neither path nor inline content)")
	}

	var cfg MultiConfig
	if err := yaml.Unmarshal([]byte(content), &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	if err := validateMultiConfig(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func validateMultiConfig(cfg *MultiConfig) error {
	if cfg.Version != "" && cfg.Version != "1" {
		return fmt.Errorf("unsupported config version %q (expected \"1\" or omit the field)", cfg.Version)
	}
	if len(cfg.Checks) == 0 {
		return fmt.Errorf("no checks defined in config")
	}
	return nil
}

// SpecToCheckContext converts a CheckSpec (plus a global default timeout) into
// a *model.CheckContext populated with the values and with CLI-matching defaults
// applied for certain fields. It does not perform validation or invoke any checker.
func SpecToCheckContext(spec CheckSpec, globalTimeout int) *model.CheckContext {
	timeout := spec.Timeout
	if timeout <= 0 {
		timeout = globalTimeout
	}
	if timeout <= 0 {
		timeout = 30
	}

	cc := &model.CheckContext{
		Host:          spec.Host,
		Port:          spec.Port,
		Username:      spec.Username,
		PasswordValue: spec.Password,
		PasswordFile:  spec.PasswordFile,
		DatabaseName:  spec.DatabaseName,
		Config: model.ConfigurationContext{
			CheckType: spec.Type,
			Timeout:   timeout,
		},
		DBConf: model.DatabaseSpecificConf{
			SSLMode:       spec.SSLMode,
			OperationMode: spec.OperationMode,
		},
		HttpConf: model.HttpSpecificConf{
			StatusCode:     spec.HttpStatusCode,
			Text:           spec.HttpText,
			FollowRedirect: true,
		},
		KafkaConf: model.KafkaConf{
			ConnectionType:              spec.KafkaConnectionType,
			ConnectToLeaderViaNonLeader: spec.KafkaConnectToLeaderViaNonLeader,
		},
		DNSConf: model.DNSConf{
			RecordType: spec.DNSType,
			Expected:   spec.DNSExpect,
			Server:     spec.DNSServer,
		},
		InfluxConf: model.InfluxConf{
			Token:  spec.InfluxToken,
			Org:    spec.InfluxOrg,
			Bucket: spec.InfluxBucket,
		},
	}

	if spec.HttpFollowRedirect != nil {
		cc.HttpConf.FollowRedirect = *spec.HttpFollowRedirect
	}

	// Apply defaults that match current CLI flag defaults (when omitted in YAML)
	if cc.DNSConf.RecordType == "" {
		cc.DNSConf.RecordType = "A"
	}
	if cc.KafkaConf.ConnectionType == "" {
		cc.KafkaConf.ConnectionType = "tcp"
	}
	if cc.HttpConf.StatusCode == 0 {
		cc.HttpConf.StatusCode = 200
	}

	return cc
}

// runMulti is the pure orchestration skeleton for multi-check mode.
// It uses the injected runner (instead of real check modules) so all
// decision logic (names, timeouts, optional, fail_fast, warnings) can be
// unit tested without Docker or real services.
func runMulti(cfg *MultiConfig, runner func(*model.CheckContext) error) error {
	if cfg == nil || len(cfg.Checks) == 0 {
		return fmt.Errorf("no checks defined")
	}

	globalTimeout := cfg.Timeout
	if globalTimeout <= 0 {
		globalTimeout = 30
	}

	failFast := true
	if cfg.FailFast != nil {
		failFast = *cfg.FailFast
	}

	var firstRequiredErr error
	var optWarnings []string

	for i, spec := range cfg.Checks {
		name := spec.Name
		if name == "" {
			name = fmt.Sprintf("check-%d", i+1)
		}
		typ := spec.Type

		display := name
		if typ != "" {
			display += fmt.Sprintf(" (%s)", typ)
		}
		if spec.Optional {
			display += ", optional"
		}
		fmt.Printf("Waiting for %s...\n", display)

		cc := SpecToCheckContext(spec, globalTimeout)

		err := runner(cc)
		if err != nil {
			if spec.Optional {
				optWarnings = append(optWarnings, fmt.Sprintf(`Warning: optional check "%s" (%s) failed: %v`, name, typ, err))
				continue
			}
			if firstRequiredErr == nil {
				firstRequiredErr = fmt.Errorf(`check "%s" (%s) failed: %w`, name, typ, err)
			}
			if failFast {
				return firstRequiredErr
			}
			// !failFast: continue running remaining checks
		}
	}

	if firstRequiredErr != nil {
		return firstRequiredErr
	}

	// overall success (possibly with optional warnings)
	for _, w := range optWarnings {
		fmt.Println(w)
	}
	if len(optWarnings) > 0 {
		fmt.Printf("Success with %d optional warning(s)!\n", len(optWarnings))
	} else {
		fmt.Println("Success!")
	}
	return nil
}

// RunMultiChecks executes the checks defined in the MultiConfig in order,
// using the real check implementations (findCheckModule + ticker).
// It reuses the decision logic from runMulti (for optional, fail_fast, warnings, etc.)
// and the core wait machinery.
func RunMultiChecks(ctx context.Context, cfg *MultiConfig) error {
	if ctx == nil {
		ctx = context.Background()
	}
	return runMulti(cfg, func(cc *model.CheckContext) error {
		return RunCheck(ctx, cc)
	})
}
