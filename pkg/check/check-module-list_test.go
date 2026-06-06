package check

import (
	"testing"

	"wait4it/pkg/model"
)

func TestCheckModuleList_AllTypesRegistered(t *testing.T) {
	expectedTypes := []string{
		"tcp",
		"mysql",
		"postgres",
		"http",
		"mongo",
		"oracle",
		"influxdb",
		"redis",
		"rabbitmq",
		"memcached",
		"elasticsearch",
		"aerospike",
		"kafka",
		"dns",
	}

	for _, checkType := range expectedTypes {
		t.Run(checkType, func(t *testing.T) {
			factory, ok := cm[checkType]
			if !ok {
				t.Errorf("check type %q not registered in module list", checkType)
				return
			}
			if factory == nil {
				t.Errorf("check type %q has nil factory", checkType)
			}
		})
	}
}

func TestCheckModuleList_NoExtraTypes(t *testing.T) {
	expectedCount := 14
	if len(cm) != expectedCount {
		t.Errorf("cm has %d entries, want %d (new check type may need to be added to this test)", len(cm), expectedCount)
	}
}

func TestCheckModuleList_FactoriesReturnCheckInterface(t *testing.T) {
	for checkType, factory := range cm {
		t.Run(checkType, func(t *testing.T) {
			// Create a minimal CheckContext that should pass validation for each type.
			// Some types require more fields than others, so we provide a generous default.
			cc := &model.CheckContext{
				Host:     "localhost",
				Port:     3306,
				Username: "testuser",
				PasswordValue: "testpass",
				DatabaseName:  "testdb",
				Config: model.ConfigurationContext{
					CheckType: checkType,
					Timeout:   1,
				},
				HttpConf: model.HttpSpecificConf{
					StatusCode: 200,
				},
				DNSConf: model.DNSConf{
					RecordType: "A",
				},
				KafkaConf: model.KafkaConf{
					ConnectionType: "tcp",
				},
				DBConf: model.DatabaseSpecificConf{
					OperationMode: "standalone",
				},
			}

			checker, err := factory(cc)
			// We don't require success since some types need real services,
			// but we verify the factory doesn't panic and returns the right types.
			if err != nil {
				// Factory returned an error (likely validation), which is acceptable
				// since we're not providing real service configs for all types.
				return
			}

			if checker == nil {
				t.Errorf("factory for %q returned nil checker without error", checkType)
			}
		})
	}
}

func TestFindCheckModule_KnownType(t *testing.T) {
	cc := &model.CheckContext{
		Host: "localhost",
		Port: 3306,
		Config: model.ConfigurationContext{
			CheckType: "tcp",
			Timeout:   1,
		},
	}

	_, err := findCheckModule(cc)
	if err != nil {
		t.Errorf("findCheckModule() for known type 'tcp' returned error: %v", err)
	}
}

func TestFindCheckModule_UnknownType(t *testing.T) {
	cc := &model.CheckContext{
		Host: "localhost",
		Port: 3306,
		Config: model.ConfigurationContext{
			CheckType: "unknown",
			Timeout:   1,
		},
	}

	_, err := findCheckModule(cc)
	if err == nil {
		t.Error("findCheckModule() should return error for unknown check type")
	}
}
