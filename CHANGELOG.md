# Changelog

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
