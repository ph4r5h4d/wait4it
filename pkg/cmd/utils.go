package cmd

import (
	"strconv"
	"time"

	"wait4it/pkg/environment"
)

func defaultEnvString(env environment.Environment, key string, fallback string) string {
	v := env.GetEnv(key)
	if v == "" {
		return fallback
	}

	return v
}

func defaultEnvInt(env environment.Environment, key string, fallback int) int {
	v := env.GetEnv(key)
	if v == "" {
		return fallback
	}

	i, err := strconv.ParseInt(v, 10, 0)
	if err != nil {
		return fallback
	}

	return int(i)
}

func defaultEnvUint(env environment.Environment, key string, fallback uint) uint {
	v := env.GetEnv(key)
	if v == "" {
		return fallback
	}

	i, err := strconv.ParseUint(v, 10, 0)
	if err != nil {
		return fallback
	}

	return uint(i)
}

func defaultEnvDuration(env environment.Environment, key string, fallback time.Duration) time.Duration {
	v := env.GetEnv(key)
	if v == "" {
		return fallback
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}

	return d
}

func defaultEnvBool(env environment.Environment, key string, fallback bool) bool {
	v := env.GetEnv(key)
	if v == "" {
		return fallback
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}

	return b
}
