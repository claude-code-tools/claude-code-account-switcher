// Package statusline implements the `claude-accounts statusline` command that
// Claude Code runs as its status line. It prints the active account name (so a
// session always shows which subscription it uses) plus live usage, and caches
// the rate limits for `doctor` / `list`. It is pure Go — no zsh dependency.
package statusline

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/leegunwoo98/claude-code-account-switcher/internal/paths"
)

type window struct {
	Used *float64 `json:"used_percentage"`
}

type payload struct {
	RateLimits *struct {
		FiveHour window `json:"five_hour"`
		SevenDay window `json:"seven_day"`
	} `json:"rate_limits"`
}

// Run reads the status-line JSON from stdin and prints the status line.
func Run() {
	label := os.Getenv("CLAUDE_SUBSCRIPTION_LABEL")
	slug := os.Getenv("CLAUDE_SUBSCRIPTION_SLUG")
	if label == "" {
		label = slug
	}
	if label == "" {
		return // not launched through the switcher
	}

	data, _ := io.ReadAll(os.Stdin)
	var p payload
	_ = json.Unmarshal(data, &p)

	if slug != "" && validSlug(slug) {
		cacheUsage(slug, data)
	}

	out := label
	if p.RateLimits != nil {
		var parts []string
		if u := p.RateLimits.FiveHour.Used; u != nil {
			parts = append(parts, fmt.Sprintf("5h %.0f%%", *u))
		}
		if u := p.RateLimits.SevenDay.Used; u != nil {
			parts = append(parts, fmt.Sprintf("7d %.0f%%", *u))
		}
		if len(parts) > 0 {
			out += "  ·  " + strings.Join(parts, " · ")
		}
	}
	fmt.Print(out)
}

// cacheUsage extracts rate_limits from the payload and writes the per-account
// usage snapshot read by doctor/list.
func cacheUsage(slug string, data []byte) {
	var raw struct {
		RateLimits json.RawMessage `json:"rate_limits"`
	}
	if json.Unmarshal(data, &raw) != nil || len(raw.RateLimits) == 0 {
		return
	}
	body, err := json.Marshal(raw)
	if err != nil {
		return
	}
	dir := paths.UsageDir()
	if os.MkdirAll(dir, 0o700) != nil {
		return
	}
	tmp := filepath.Join(dir, slug+".json.tmp")
	if os.WriteFile(tmp, body, 0o600) == nil {
		_ = os.Rename(tmp, filepath.Join(dir, slug+".json"))
	}
}

func validSlug(s string) bool {
	for _, r := range s {
		if !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-') {
			return false
		}
	}
	return s != ""
}
