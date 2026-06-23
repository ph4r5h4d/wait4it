package main

import (
	"testing"
)

func TestDefaultEnvBool(t *testing.T) {
	tests := []struct {
		name string
		env  string
		val  string
		def  bool
		want bool
	}{
		{"unset uses default", "W4IT_TEST_BOOL_UNSET", "", false, false},
		{"default true", "W4IT_TEST_BOOL_UNSET", "", true, true},
		{"true string", "W4IT_TEST_BOOL_TRUE", "true", false, true},
		{"one string", "W4IT_TEST_BOOL_ONE", "1", false, true},
		{"false string", "W4IT_TEST_BOOL_FALSE", "false", true, false},
		{"zero string", "W4IT_TEST_BOOL_ZERO", "0", true, false},
		{"garbage returns false", "W4IT_TEST_BOOL_GARBAGE", "garbage", false, false},
		{"garbage returns false ignoring default", "W4IT_TEST_BOOL_GARBAGE", "garbage", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.val != "" {
				t.Setenv(tt.env, tt.val)
			}
			got := defaultEnvBool(tt.env, tt.def)
			if got != tt.want {
				t.Fatalf("defaultEnvBool(%q, %v) = %v, want %v", tt.env, tt.def, got, tt.want)
			}
		})
	}
}