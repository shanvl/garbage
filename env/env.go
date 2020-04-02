// Package env helps parse environment variables
package env

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Int parses env as an integer. If env is empty or not a valid int, Int returns a fallback
func Int(name string, fallback int) int {
	env := os.Getenv(name)
	if env == "" {
		return fallback
	}
	envInt, err := strconv.Atoi(env)
	if err != nil {
		fmt.Println(envInt)
		return fallback
	}
	return envInt
}

// Minutes parses env as minutes. If env is empty or not a valid time.Duration, Minutes returns a fallback
func Minutes(name string, fallback time.Duration) time.Duration {
	return time.Minute * duration(name, fallback)
}

// Seconds parses env as minutes. If env is empty or not a valid time.Duration, Seconds returns a fallback
func Seconds(name string, fallback time.Duration) time.Duration {
	return time.Second * duration(name, fallback)
}

// String parses env as a string. If env is empty, String returns a fallback
func String(name string, fallback string) string {
	env := os.Getenv(name)
	if env == "" {
		return fallback
	}
	return env
}

// duration parses env as time.Duration. If env is empty or not a valid time.Duration, duration returns a fallback
func duration(name string, fallback time.Duration) time.Duration {
	env := os.Getenv(name)
	if env == "" {
		return fallback
	}
	envInt, err := strconv.Atoi(env)
	if err != nil {
		return fallback
	}
	return time.Duration(envInt)
}
