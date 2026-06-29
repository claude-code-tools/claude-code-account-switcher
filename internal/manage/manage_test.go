package manage

import "testing"

func TestSlugify(t *testing.T) {
	for in, want := range map[string]string{
		"Gmail Work":   "gmail-work",
		"claude-naver": "naver",
		"  Spaces  ":   "spaces",
		"A--B":         "a-b",
		"Work_2!":      "work2",
		"Personal":     "personal",
	} {
		if got := Slugify(in); got != want {
			t.Errorf("Slugify(%q) = %q, want %q", in, got, want)
		}
	}
}
