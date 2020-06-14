// Package env helps parse environment variables
package env

import (
	"os"
	"strconv"
	"time"
)

// Bool parses env as an bool. If env is empty or not a valid bool, Bool returns a fallback
func Bool(name string, fallback bool) bool {
	env := os.Getenv(name)
	if env == "" {
		return fallback
	}
	envBool, err := strconv.ParseBool(env)
	if err != nil {
		return fallback
	}
	return envBool
}

// Int parses env as an integer. If env is empty or not a valid int, Int returns a fallback
func Int(name string, fallback int) int {
	env := os.Getenv(name)
	if env == "" {
		return fallback
	}
	envInt, err := strconv.Atoi(env)
	if err != nil {
		return fallback
	}
	return envInt
}

// Duration parses env as duration. If env is empty or not a valid time.Duration, Duration returns a fallback
func Duration(name string, fallback time.Duration) time.Duration {
	env := os.Getenv(name)
	if env == "" {
		return fallback
	}
	duration, err := time.ParseDuration(env)
	if err != nil {
		return fallback
	}
	return duration
}

// String parses env as a string. If env is empty, String returns a fallback
func String(name string, fallback string) string {
	env := os.Getenv(name)
	if env == "" {
		return fallback
	}
	return env
}
