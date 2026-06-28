# Changelog

## 0.3.0 - 2026-06-28

- Fix the core bug where every subscription billed the account from `oauth/login`
  regardless of the injected token. Claude Code pins the account/organization
  cached in `~/.claude.json`'s `oauthAccount` to every session, so injecting
  `CLAUDE_CODE_OAUTH_TOKEN` alone still resolved to the logged-in plan. Each
  `claude-<suffix>` now launches with its own `CLAUDE_CONFIG_DIR` that strips
  only `oauthAccount` and symlinks everything else, so the right plan is billed
  while settings, plugins, memory, and history stay shared. Requires `jq` or
  `python3`; without either, the session launches unisolated with a warning.
- Show the pinned account name on the status line for the whole session, so an
  in-session `/login` can't silently mislead which subscription is active.
- Keep the test suite hermetic against leaked `CLAUDE_ACCOUNTS_*` environment
  variables so a developer's shell can't make the installer test escape its
  sandbox.

## 0.2.4 - 2026-06-28

- Remove the optional per-email `claude auth login` step from token generation.
  It wrote global subscription credentials that the injected
  `CLAUDE_CODE_OAUTH_TOKEN` always overrides, so it had no effect on which
  account a session used, and its unconditional `auth logout` could destroy the
  active Claude CLI session if the login failed. `claude setup-token` already
  binds the token to whichever account is signed in on claude.ai; the prompt now
  says so.
- Clear cached usage only after a new profile is registered, so a failed save no
  longer discards usage for an account that was never added.

## 0.2.3 - 2026-06-28

- Optionally force a full Claude login for a specific email before generating or
  refreshing a setup token.
- Clear cached usage whenever a profile token is replaced.

## 0.2.2 - 2026-06-28

- Pin installer downloads to one release so raw GitHub caching cannot mix files
  from different commits.
- Cache-bust explicitly requested mutable versions such as `main`.

## 0.2.1 - 2026-06-28

- Show the selected profile name in Claude's session header, terminal title,
  and resume picker by default.

## 0.2.0 - 2026-06-27

- Install real `claude-accounts` and `claude-<suffix>` executables on `PATH`.
- Stop modifying shell startup files during installation.

## 0.1.1 - 2026-06-27

- Add a reproducible animated terminal showcase to the README.

## 0.1.0 - 2026-06-27

- Add Keychain-backed Claude Code subscription profiles.
- Add generated `claude-<suffix>` commands.
- Add unified `claude-accounts` management UI.
- Add cached 5-hour and 7-day usage display.
- Add installer and uninstaller for Zsh on macOS.
