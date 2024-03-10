package utils

import (
	"os"
	"testing"
)

// LoadEnv loads the environment variable and returns it for testing. If the environment variable is not set, it skips the test.
func LoadEnv(t *testing.T, key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		t.Skipf("Skip test because %s is not set", key)
	}

	return v
}
