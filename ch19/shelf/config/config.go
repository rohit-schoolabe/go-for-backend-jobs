// Package config reads Shelf's settings from the
// environment, the one place a container can set them.
package config

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Config is every setting the server reads at start.
type Config struct {
	Addr            string        // ADDR
	DatabaseURL     string        // DATABASE_URL
	JWTSecret       string        // JWT_SECRET
	LoanDays        int           // LOAN_DAYS
	DrainDelay      time.Duration // DRAIN_DELAY
	ShutdownTimeout time.Duration // SHUTDOWN_TIMEOUT
}

// Load builds a Config from getenv, which is os.Getenv in
// main and a map lookup in tests. It reports every bad
// value at once.
func Load(getenv func(string) string) (Config, error) {
	c := Config{
		Addr:            or(getenv("ADDR"), ":8080"),
		DatabaseURL:     or(getenv("DATABASE_URL"), "shelf.db"),
		JWTSecret:       getenv("JWT_SECRET"),
		LoanDays:        14,
		DrainDelay:      5 * time.Second,
		ShutdownTimeout: 20 * time.Second,
	}
	var errs []error
	if len(c.JWTSecret) < 32 {
		errs = append(errs, errors.New(
			"JWT_SECRET must be 32 or more characters"))
	}
	if v := getenv("LOAN_DAYS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			errs = append(errs, fmt.Errorf(
				"LOAN_DAYS=%q: want a positive number", v))
		}
		c.LoanDays = n
	}
	durations := []struct {
		key string
		dst *time.Duration
	}{
		{"DRAIN_DELAY", &c.DrainDelay},
		{"SHUTDOWN_TIMEOUT", &c.ShutdownTimeout},
	}
	for _, d := range durations {
		v := getenv(d.key)
		if v == "" {
			continue
		}
		parsed, err := time.ParseDuration(v)
		if err != nil || parsed < 0 {
			errs = append(errs, fmt.Errorf(
				"%s=%q: want a duration like 10s", d.key, v))
		}
		*d.dst = parsed
	}
	return c, errors.Join(errs...)
}

func or(v, def string) string {
	if v == "" {
		return def
	}
	return v
}
