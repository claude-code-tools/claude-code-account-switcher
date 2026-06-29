// Command claude-accounts manages and launches Keychain-backed Claude Code
// subscriptions. When invoked through a claude-<slug> name it launches that
// subscription directly.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/leegunwoo98/claude-code-account-switcher/internal/doctor"
	"github.com/leegunwoo98/claude-code-account-switcher/internal/launcher"
	"github.com/leegunwoo98/claude-code-account-switcher/internal/registry"
	"github.com/leegunwoo98/claude-code-account-switcher/internal/usage"
)

func main() {
	prog := filepath.Base(os.Args[0])
	args := os.Args[1:]

	// Direct dispatch: invoked as claude-<slug>.
	if strings.HasPrefix(prog, "claude-") && prog != "claude-accounts" {
		launch(strings.TrimPrefix(prog, "claude-"), args)
		return
	}

	if len(args) > 0 {
		switch args[0] {
		case "doctor":
			if err := doctor.Run(); err != nil {
				fail(err)
			}
		case "list", "ls":
			list()
		case "launch", "run":
			if len(args) < 2 {
				fail(fmt.Errorf("usage: claude-accounts launch <slug> [claude args...]"))
			}
			launch(args[1], args[2:])
		case "-h", "--help", "help":
			printHelp()
		default:
			fail(fmt.Errorf("unknown command %q (try --help)", args[0]))
		}
		return
	}

	// No subcommand: list. (The interactive manager is a later milestone.)
	list()
}

func launch(slug string, args []string) {
	acct, ok := registry.Find(slug)
	if !ok {
		fail(fmt.Errorf("unknown subscription: %s\nRun: claude-accounts", slug))
	}
	if err := launcher.Launch(acct, args); err != nil {
		os.Exit(1)
	}
}

func list() {
	accounts, err := registry.Load()
	if err != nil {
		fail(err)
	}
	if len(accounts) == 0 {
		fmt.Println("No Claude subscriptions configured.")
		return
	}
	for _, a := range accounts {
		fmt.Printf("%-18s %-24s %s\n", a.Command(), a.Label, usage.Summary(a.Slug))
	}
}

func printHelp() {
	fmt.Print(`usage: claude-accounts [command]

  (no command)    list configured subscriptions
  list            list configured subscriptions
  launch <slug>   launch Claude as that subscription
  doctor          check tokens, isolation, and same-account collisions

Each configured subscription also dispatches as claude-<slug>.
`)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
