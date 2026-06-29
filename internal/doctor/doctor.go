// Package doctor diagnoses the configured subscriptions: token presence, and
// whether any two accounts resolve to the same subscription (identical token or
// identical usage fingerprint) — the symptom of a token generated under the
// wrong account. It reads only local data.
package doctor

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/leegunwoo98/claude-code-account-switcher/internal/credstore"
	"github.com/leegunwoo98/claude-code-account-switcher/internal/registry"
	"github.com/leegunwoo98/claude-code-account-switcher/internal/usage"
)

// Run prints the report and returns an error only on a registry read failure.
func Run() error {
	accounts, err := registry.Load()
	if err != nil {
		return err
	}
	if len(accounts) == 0 {
		fmt.Println("No Claude subscriptions configured. Run: claude-accounts")
		return nil
	}

	fmt.Println("Claude subscription doctor")
	fmt.Print("  per-account config isolation: native (no external dependency)\n\n")

	tokenOwner := map[string]string{}
	usageOwner := map[string]string{}
	var warnings []string

	for _, a := range accounts {
		fmt.Printf("● %s  (%s)\n", a.Label, a.Command())

		token, _ := credstore.Get(a.Service)
		if token == "" {
			fmt.Println("    token : MISSING — refresh it from claude-accounts")
			warnings = append(warnings, a.Label+": no Keychain token")
			fmt.Println()
			continue
		}
		if strings.HasPrefix(token, "sk-ant-oat") {
			fmt.Println("    token : present (setup-token)")
		} else {
			fmt.Println("    token : present, but unexpected prefix")
			warnings = append(warnings, a.Label+": token does not look like a setup-token")
		}

		sum := sha256.Sum256([]byte(token))
		h := hex.EncodeToString(sum[:])
		if owner, ok := tokenOwner[h]; ok {
			fmt.Printf("    ⚠ identical token to %q (same account billed for both)\n", owner)
			warnings = append(warnings, fmt.Sprintf("%s and %s share one token", a.Label, owner))
		} else {
			tokenOwner[h] = a.Label
		}

		fmt.Printf("    usage : %s\n", usage.Summary(a.Slug))
		if fp := usage.Fingerprint(a.Slug); fp != "" {
			if owner, ok := usageOwner[fp]; ok {
				fmt.Printf("    ⚠ identical usage to %q — likely the SAME subscription\n", owner)
				warnings = append(warnings, fmt.Sprintf(
					"%s and %s report identical usage; one token was probably generated under the wrong account",
					a.Label, owner))
			} else {
				usageOwner[fp] = a.Label
			}
		}
		fmt.Println()
	}

	if len(warnings) > 0 {
		fmt.Println("Findings:")
		for _, w := range warnings {
			fmt.Printf("  ⚠ %s\n", w)
		}
	} else {
		fmt.Println("✓ No problems detected.")
	}
	fmt.Println()
	fmt.Println("Usage reflects each account's last launch — launch one, then re-run")
	fmt.Println("'claude-accounts doctor'. Two accounts with identical usage are billing")
	fmt.Println("the same subscription.")
	return nil
}
