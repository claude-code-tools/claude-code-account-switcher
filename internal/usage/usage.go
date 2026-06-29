// Package usage reads the cached per-account rate-limit snapshots written by the
// status-line hook, and derives display summaries and same-account fingerprints.
package usage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/leegunwoo98/claude-code-account-switcher/internal/paths"
)

type window struct {
	UsedPercentage *float64 `json:"used_percentage"`
	ResetsAt       *int64   `json:"resets_at"`
}

type snapshot struct {
	RateLimits struct {
		FiveHour window `json:"five_hour"`
		SevenDay window `json:"seven_day"`
	} `json:"rate_limits"`
}

func load(slug string) (*snapshot, bool) {
	b, err := os.ReadFile(filepath.Join(paths.UsageDir(), slug+".json"))
	if err != nil {
		return nil, false
	}
	var s snapshot
	if json.Unmarshal(b, &s) != nil {
		return nil, false
	}
	return &s, true
}

// Summary renders a one-line "5h X% · 7d Y% used" string, or "usage pending".
func Summary(slug string) string {
	s, ok := load(slug)
	if !ok {
		return "usage pending"
	}
	out := ""
	if p := s.RateLimits.FiveHour.UsedPercentage; p != nil {
		out = fmt.Sprintf("5h %.0f%%", *p)
	}
	if p := s.RateLimits.SevenDay.UsedPercentage; p != nil {
		if out != "" {
			out += " · "
		}
		out += fmt.Sprintf("7d %.0f%%", *p)
	}
	if out == "" {
		return "usage pending"
	}
	return out + " used"
}

// Fingerprint is a stable signature of an account's usage, used to detect when
// two accounts are billing the same subscription. Empty when no data.
func Fingerprint(slug string) string {
	s, ok := load(slug)
	if !ok {
		return ""
	}
	f, d := s.RateLimits.FiveHour, s.RateLimits.SevenDay
	if f.UsedPercentage == nil && d.UsedPercentage == nil {
		return ""
	}
	return fstr(f.UsedPercentage) + "|" + istr(f.ResetsAt) + "|" + fstr(d.UsedPercentage) + "|" + istr(d.ResetsAt)
}

// FiveHourNearLimit reports the 5-hour usage when its window is still active and
// at or above threshold; otherwise ok is false.
func FiveHourNearLimit(slug string, threshold float64, now int64) (pct float64, ok bool) {
	s, found := load(slug)
	if !found {
		return 0, false
	}
	f := s.RateLimits.FiveHour
	if f.UsedPercentage == nil || f.ResetsAt == nil {
		return 0, false
	}
	if *f.ResetsAt <= now || *f.UsedPercentage < threshold {
		return 0, false
	}
	return *f.UsedPercentage, true
}

func fstr(p *float64) string {
	if p == nil {
		return ""
	}
	return strconv.FormatFloat(*p, 'g', -1, 64)
}

func istr(p *int64) string {
	if p == nil {
		return ""
	}
	return strconv.FormatInt(*p, 10)
}
