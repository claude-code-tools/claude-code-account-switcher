package registry

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSkipsInvalidAndDuplicate(t *testing.T) {
	dir := t.TempDir()
	tsv := filepath.Join(dir, "accounts.tsv")
	content := "" +
		"# comment\n" +
		"gmail\tGmail Work\tClaude Code Subscription: claude-gmail\n" +
		"naver\tNaver\tClaude Code Subscription: claude-naver\n" +
		"gmail\tDuplicate\tsvc\n" + // duplicate slug -> skipped
		"Bad_Slug\tx\tsvc\n" + // invalid slug -> skipped
		"missingcols\tonly-two\n" // too few columns -> skipped
	if err := os.WriteFile(tsv, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CLAUDE_SUBSCRIPTIONS_FILE", tsv)

	accounts, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(accounts) != 2 {
		t.Fatalf("got %d accounts, want 2: %+v", len(accounts), accounts)
	}
	if accounts[0].Slug != "gmail" || accounts[0].Command() != "claude-gmail" {
		t.Errorf("unexpected first account: %+v", accounts[0])
	}
	if accounts[1].Service != "Claude Code Subscription: claude-naver" {
		t.Errorf("unexpected service: %q", accounts[1].Service)
	}
}

func TestValidSlug(t *testing.T) {
	for slug, want := range map[string]bool{
		"gmail": true, "work-2": true,
		"": false, "-x": false, "x-": false, "Up": false, "a_b": false, "a b": false,
	} {
		if got := ValidSlug(slug); got != want {
			t.Errorf("ValidSlug(%q) = %v, want %v", slug, got, want)
		}
	}
}
