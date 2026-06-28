# Changelog

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
