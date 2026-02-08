package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveLoadConfigRoundTrip(t *testing.T) {
	tmp := t.TempDir()
	setUpEnvironment(t, tmp)

	cfg := &AuthConfig{
		AuthMethod:  AuthMethodAppToken,
		AccessToken: "abc123",
	}

	err := SaveConfig(cfg)
	assert.NoError(t, err)

	got, err := LoadConfig()
	assert.NoError(t, err)

	assert.Equal(t, cfg.AuthMethod, got.AuthMethod)
	assert.Equal(t, cfg.AccessToken, got.AccessToken)
}

func TestGetConfigPath(t *testing.T) {
	tmp := t.TempDir()
	setUpEnvironment(t, tmp)

	path, err := GetConfigPath()
	assert.NoError(t, err)

	var want string
	if runtime.GOOS == "darwin" {
		want = filepath.Join(tmp, "Library", "Application Support", "vkscape", "config.json")
	} else {
		want = filepath.Join(tmp, "vkscape", "config.json")
	}

	assert.Equal(t, want, path)
}

func setUpEnvironment(t *testing.T, tmp string) {
	t.Helper()

	restoreEnv(t, "HOME")
	restoreEnv(t, "XDG_CONFIG_HOME")

	if runtime.GOOS == "darwin" {
		if err := os.Setenv("HOME", tmp); err != nil {
			t.Fatalf("set HOME: %v", err)
		}
	} else {
		if err := os.Setenv("XDG_CONFIG_HOME", tmp); err != nil {
			t.Fatalf("set XDG_CONFIG_HOME: %v", err)
		}
	}
}

func restoreEnv(t *testing.T, key string) {
	t.Helper()
	val, ok := os.LookupEnv(key)
	if !ok {
		t.Cleanup(func() { _ = os.Unsetenv(key) })
		return
	}
	t.Cleanup(func() { _ = os.Setenv(key, val) })
}
