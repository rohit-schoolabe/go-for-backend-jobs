package config

import (
	"strings"
	"testing"
	"time"
)

func env(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func TestLoadDefaults(t *testing.T) {
	c, err := Load(env(map[string]string{
		"JWT_SECRET": strings.Repeat("k", 32),
	}))
	if err != nil {
		t.Fatal(err)
	}
	if c.Addr != ":8080" || c.LoanDays != 14 ||
		c.ShutdownTimeout != 20*time.Second {
		t.Errorf("defaults wrong: %+v", c)
	}
}

func TestLoadReportsEveryProblem(t *testing.T) {
	_, err := Load(env(map[string]string{
		"LOAN_DAYS":        "two weeks",
		"SHUTDOWN_TIMEOUT": "20",
	}))
	if err == nil {
		t.Fatal("want an error")
	}
	for _, key := range []string{
		"JWT_SECRET", "LOAN_DAYS", "SHUTDOWN_TIMEOUT",
	} {
		if !strings.Contains(err.Error(), key) {
			t.Errorf("error does not mention %s: %v",
				key, err)
		}
	}
}
