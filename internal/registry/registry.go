// Package registry reads the tab-separated account list (slug, label, Keychain
// service). It never stores or returns tokens — those live in the OS credential
// store, keyed by the service name recorded here.
package registry

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/leegunwoo98/claude-code-account-switcher/internal/paths"
)

// Account is one configured subscription.
type Account struct {
	Slug    string
	Label   string
	Service string
}

// Command is the direct launcher name for this account, e.g. "claude-gmail".
func (a Account) Command() string { return "claude-" + a.Slug }

// Load returns the configured accounts, skipping malformed or duplicate rows.
func Load() ([]Account, error) {
	f, err := os.Open(paths.AccountsFile())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var accounts []Account
	seen := map[string]bool{}
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "\t", 4)
		if len(parts) < 3 {
			continue
		}
		slug, label, service := parts[0], parts[1], parts[2]
		if !ValidSlug(slug) || label == "" || service == "" || seen[slug] {
			continue
		}
		seen[slug] = true
		accounts = append(accounts, Account{Slug: slug, Label: label, Service: service})
	}
	return accounts, sc.Err()
}

// Find returns the account for a slug.
func Find(slug string) (Account, bool) {
	accounts, _ := Load()
	for _, a := range accounts {
		if a.Slug == slug {
			return a, true
		}
	}
	return Account{}, false
}

// ServiceFor is the Keychain service name used for a slug's token.
func ServiceFor(slug string) string {
	return fmt.Sprintf("Claude Code Subscription: claude-%s", slug)
}

// Append adds an account row to the registry, creating the file if needed.
func Append(a Account) error {
	if err := os.MkdirAll(filepath.Dir(paths.AccountsFile()), 0o700); err != nil {
		return err
	}
	f, err := os.OpenFile(paths.AccountsFile(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%s\t%s\t%s\n", a.Slug, a.Label, a.Service)
	return err
}

// Remove rewrites the registry without the given slug.
func Remove(slug string) error {
	accounts, err := Load()
	if err != nil {
		return err
	}
	var b strings.Builder
	for _, a := range accounts {
		if a.Slug == slug {
			continue
		}
		fmt.Fprintf(&b, "%s\t%s\t%s\n", a.Slug, a.Label, a.Service)
	}
	tmp := paths.AccountsFile() + ".tmp"
	if err := os.WriteFile(tmp, []byte(b.String()), 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, paths.AccountsFile())
}

// ValidSlug matches the zsh validator: lowercase [a-z0-9-], no leading/trailing
// hyphen, non-empty.
func ValidSlug(s string) bool {
	if s == "" || strings.HasPrefix(s, "-") || strings.HasSuffix(s, "-") {
		return false
	}
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-':
		default:
			return false
		}
	}
	return true
}
