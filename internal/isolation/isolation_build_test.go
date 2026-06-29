package isolation

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestBuildSharesAndStrips exercises the full per-account config build on every
// platform: shared entries are reachable through the platform link (symlink on
// Unix; symlink/junction/hardlink/copy on Windows), the account is stripped, and
// per-process runtime is not shared.
func TestBuildSharesAndStrips(t *testing.T) {
	base := t.TempDir()
	t.Setenv("CLAUDE_CONFIG_DIR", base)
	t.Setenv("CLAUDE_SUBSCRIPTIONS_CONFIG_DIR", t.TempDir())

	mustWriteFile(t, filepath.Join(base, ".claude.json"), `{"oauthAccount":{"org":"x"},"keep":0.15}`)
	mustWriteFile(t, filepath.Join(base, "settings.json"), `{"model":"opus"}`)
	if err := os.MkdirAll(filepath.Join(base, "plugins"), 0o700); err != nil {
		t.Fatal(err)
	}
	mustWriteFile(t, filepath.Join(base, "plugins", "p.txt"), "plugin")
	mustWriteFile(t, filepath.Join(base, "daemon.lock"), "lock")

	dir, err := Build("acct")
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(mustReadFile(t, filepath.Join(dir, ".claude.json"))), &m); err != nil {
		t.Fatalf("per-account config is not valid JSON: %v", err)
	}
	if _, ok := m["oauthAccount"]; ok {
		t.Error("oauthAccount was not stripped")
	}
	if string(m["keep"]) != "0.15" {
		t.Errorf("keep = %s, want 0.15 (exact)", m["keep"])
	}

	if got := mustReadFile(t, filepath.Join(dir, "settings.json")); got != `{"model":"opus"}` {
		t.Errorf("shared settings.json = %q", got)
	}
	if got := mustReadFile(t, filepath.Join(dir, "plugins", "p.txt")); got != "plugin" {
		t.Errorf("shared plugins/p.txt = %q", got)
	}

	if _, err := os.Lstat(filepath.Join(dir, "daemon.lock")); err == nil {
		t.Error("daemon.lock should not be shared into the account dir")
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func mustReadFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
