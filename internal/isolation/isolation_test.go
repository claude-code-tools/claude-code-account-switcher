package isolation

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteStrippedRemovesAccountPreservesRest(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.json")
	dst := filepath.Join(dir, "dst.json")

	// Includes a float, a null, and nesting — all must survive untouched.
	input := `{"oauthAccount":{"organizationUuid":"x"},"numStartups":7,"ratio":0.15,"nilable":null,"nested":{"a":[1,2,3]}}`
	if err := os.WriteFile(src, []byte(input), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := writeStripped(src, dst); err != nil {
		t.Fatal(err)
	}

	b, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]json.RawMessage
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if _, ok := got["oauthAccount"]; ok {
		t.Error("oauthAccount was not removed")
	}
	for key, want := range map[string]string{
		"numStartups": "7",
		"ratio":       "0.15", // exact: RawMessage avoids float re-serialization
		"nilable":     "null",
		"nested":      `{"a":[1,2,3]}`,
	} {
		if got := string(got[key]); got != want {
			t.Errorf("%s = %s, want %s", key, got, want)
		}
	}
}

func TestWriteStrippedMissingSource(t *testing.T) {
	dir := t.TempDir()
	dst := filepath.Join(dir, "dst.json")
	if err := writeStripped(filepath.Join(dir, "absent.json"), dst); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(dst)
	if string(b) != "{}\n" {
		t.Errorf("missing-source output = %q, want %q", string(b), "{}\n")
	}
}

func TestIsolatedName(t *testing.T) {
	cases := map[string]bool{
		"settings.json": false,
		"plugins":       false,
		"daemon":        true,
		"daemon.lock":   true,
		"foo.lock":      true,
		"bridge.sock":   true,
	}
	for name, want := range cases {
		if got := isolatedName(name); got != want {
			t.Errorf("isolatedName(%q) = %v, want %v", name, got, want)
		}
	}
}
