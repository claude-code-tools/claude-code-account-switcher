// Package paths resolves the on-disk locations the switcher uses, honoring the
// same environment overrides as the zsh implementation so the two stay
// interoperable on macOS.
package paths

import (
	"os"
	"path/filepath"
)

func home() string {
	if h, err := os.UserHomeDir(); err == nil && h != "" {
		return h
	}
	return os.Getenv("HOME")
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// ConfigRoot is the switcher's own config directory.
func ConfigRoot() string {
	if v := os.Getenv("CLAUDE_SUBSCRIPTIONS_DIR"); v != "" {
		return v
	}
	base := envOr("XDG_CONFIG_HOME", filepath.Join(home(), ".config"))
	return filepath.Join(base, "claude-subscriptions")
}

// AccountsFile is the tab-separated registry of accounts (never holds tokens).
func AccountsFile() string {
	return envOr("CLAUDE_SUBSCRIPTIONS_FILE", filepath.Join(ConfigRoot(), "accounts.tsv"))
}

// UsageDir holds cached per-account rate-limit snapshots.
func UsageDir() string {
	return envOr("CLAUDE_SUBSCRIPTIONS_USAGE_DIR", filepath.Join(ConfigRoot(), "usage"))
}

// UsageSettings is the Claude --settings file that registers the status line.
func UsageSettings() string {
	return envOr("CLAUDE_SUBSCRIPTIONS_USAGE_SETTINGS", filepath.Join(ConfigRoot(), "usage-settings.json"))
}

// ConfigsDir is the root for per-account CLAUDE_CONFIG_DIRs.
func ConfigsDir() string {
	return envOr("CLAUDE_SUBSCRIPTIONS_CONFIG_DIR", filepath.Join(ConfigRoot(), "configs"))
}

// ClaudeHome is the base Claude config directory whose contents are shared into
// each per-account config dir.
func ClaudeHome() string {
	if v := os.Getenv("CLAUDE_CONFIG_DIR"); v != "" {
		return v
	}
	return filepath.Join(home(), ".claude")
}

// ClaudeJSON is the base .claude.json (the file that caches oauthAccount).
func ClaudeJSON() string {
	if v := os.Getenv("CLAUDE_CONFIG_DIR"); v != "" {
		return filepath.Join(v, ".claude.json")
	}
	return filepath.Join(home(), ".claude.json")
}
