# Claude Code Account Switcher

Run multiple Claude Code subscriptions side by side — on macOS, Linux, and Windows — while sharing the same `~/.claude` settings, plugins, skills, and session history.

The switcher stores one long-lived OAuth token per account in your OS credential store and gives each account a direct command such as `claude-gmail` and `claude-naver`. It also provides an account manager for adding, launching, editing, refreshing, and removing accounts. It ships as a single self-contained binary — no shell, runtime, or interpreter required.

![Claude Code Account Switcher terminal demo](assets/demo.gif)

## Why this exists

Running two Claude subscriptions on one machine is harder than it looks. Claude Code stores `/login` credentials in a single shared entry, so logging into another subscription changes the account used by every other CLI session. `claude setup-token` mints a one-year OAuth token per account, but prints it without saving it.

The subtle part: injecting `CLAUDE_CODE_OAUTH_TOKEN` alone is **not enough**. Claude Code caches the signed-in account and organization in `~/.claude.json` (`oauthAccount`) and pins that organization to every request — so a session launched with account B's token still bills account A (a `403 "missing scope or org access"` followed by a silent fallback to the logged-in account). Most multi-account setups dodge this by giving each account a fully separate `CLAUDE_CONFIG_DIR`, which works but silos your settings, MCP servers, memory, and session history per account.

This project does both halves: it stores each token under its own credential-store entry and injects only the selected one, **and** it gives each account a `CLAUDE_CONFIG_DIR` that strips only the cached `oauthAccount` while symlinking everything else back from `~/.claude`. The result is one shared configuration with the account — and its billing — correctly isolated per command.

## Requirements

- macOS, Linux, or Windows
- [Claude Code](https://code.claude.com/docs/en/setup)
- `curl` (for the install script)
- Optional: [`fzf`](https://github.com/junegunn/fzf) for the full-screen account selector (a numbered menu is used when it is absent)

Account isolation and usage parsing are built in — no `jq`, `python3`, or other external tools are needed.

## Install

### macOS and Linux

Review the installer before running it:

```sh
curl -fsSL https://raw.githubusercontent.com/leegunwoo98/claude-code-account-switcher/v1.0.0/install.sh
```

Then install:

```sh
curl -fsSL https://raw.githubusercontent.com/leegunwoo98/claude-code-account-switcher/v1.0.0/install.sh | sh
```

The installer downloads the prebuilt `claude-accounts` binary for your OS and architecture, installs it into a writable directory already on your `PATH` (defaulting to `~/.local/bin`), and does **not** modify `.zshrc`, `.bashrc`, or any other shell startup file. Set `CLAUDE_ACCOUNTS_BIN_DIR` to choose the install directory.

### Windows

Download `claude-accounts_v1.0.0_windows_amd64.zip` (or `_arm64`) from the [latest release](https://github.com/leegunwoo98/claude-code-account-switcher/releases), extract `claude-accounts.exe`, and place it in a directory on your `PATH`.

### From source (any platform)

```sh
go install github.com/leegunwoo98/claude-code-account-switcher/cmd/claude-accounts@latest
```

After installing, run the manager once to register an account; each account you add also creates its own `claude-<suffix>` launcher next to the binary:

```sh
claude-accounts
```

## Usage

Open the account manager:

```sh
claude-accounts
```

From the UI you can:

- Add a subscription
- Launch it
- Edit its display name or direct command
- Replace its token
- Remove it

During **Add** or **Refresh**, the switcher runs `claude setup-token` with all competing credentials scrubbed from the environment, so the token binds to whichever account is signed in on claude.ai in your browser — not to an already-injected token. Switch to the intended account on claude.ai first, then confirm it in the tab that opens and paste the printed token. Replacing a token also clears that account's cached usage so stale values are not attributed to the new token.

When adding an account, choose a display name and command suffix. If the suffix is `gmail`, the generated command is:

```sh
claude-gmail
```

All Claude arguments are forwarded:

```sh
claude-gmail --continue
claude-naver --model opus
claude-gmail "Review this repository"
```

The active account is shown in the session **status line** for the whole session, so an in-session `/login` can't silently mislead you about which subscription is billing. This keeps Claude's own auto-generated descriptions in the resume picker intact. To instead name the session after the account (it then appears in Claude's header, terminal title, and resume picker, replacing the auto description), set `CLAUDE_SUBSCRIPTION_NAME_SESSIONS=1`. Passing your own `--name`/`-n` always wins.

## Usage percentages

Claude Code exposes subscription rate-limit data to status-line commands after a successful API response. Wrapper-launched sessions cache those values locally without making extra requests.

The account manager displays cached values like:

```text
5h 24% · 7d 41% used
```

New accounts show `usage pending` until their first normal Claude response. Cached values may be stale until that account is used again.

## Diagnostics

```sh
claude-accounts doctor
```

`doctor` checks each account's stored token and reports when two accounts resolve to the **same** subscription — either an identical token or an identical usage fingerprint — which is the symptom of a token generated under the wrong account. Launch an account first to refresh its cached usage, then re-run.

## Token handling and security

- Tokens are generated by the official `claude setup-token` command and are valid for one year.
- Claude intentionally prints setup tokens without saving them. Copy the token once and paste it into the switcher's prompt.
- Tokens are stored in your OS credential store, never in the account registry: the **macOS Keychain**, or a `0600`-permission `tokens.json` under `~/.config/claude-subscriptions` on Linux and Windows (the same file model Claude Code uses, which works headless on servers and in containers).
- `CLAUDE_CODE_SUBPROCESS_ENV_SCRUB=1` and unset `ANTHROPIC_*` / Bedrock / Vertex / Foundry variables prevent Bash tools, hooks, and stdio MCP servers from inheriting Anthropic credentials.
- A token included in chat, logs, screenshots, or shell history must be revoked and replaced immediately.
- Replacing a locally stored token does not necessarily revoke the old server-side token. Revoke the old token in Claude settings.
- Do not run `/login` inside a wrapper-launched session; it changes authentication behavior for that session.

## Files

| Path | Purpose |
| --- | --- |
| A writable directory on `PATH` | `claude-accounts` binary and per-account `claude-<suffix>` launchers |
| `~/.config/claude-subscriptions/accounts.tsv` | Display names, command suffixes, and credential-store keys; no tokens |
| `~/.config/claude-subscriptions/configs/<suffix>/` | Per-account `CLAUDE_CONFIG_DIR` (shared symlinks + an `oauthAccount`-stripped `.claude.json`) |
| `~/.config/claude-subscriptions/usage/` | Cached rate-limit percentages |
| macOS Keychain, or `~/.config/claude-subscriptions/tokens.json` (Linux/Windows) | OAuth tokens |

## Update

Re-run the installer (macOS/Linux), download the newer release (Windows), or `go install ...@latest`. Account metadata and stored tokens are preserved.

```sh
curl -fsSL https://raw.githubusercontent.com/leegunwoo98/claude-code-account-switcher/v1.0.0/install.sh | sh
```

## Uninstall

Remove the binary and its generated launchers from your install directory:

```sh
rm -f ~/.local/bin/claude-accounts ~/.local/bin/claude-*
```

To also remove the registry, per-account config dirs, and usage cache:

```sh
rm -rf ~/.config/claude-subscriptions
```

On macOS, delete each `Claude Code Subscription: claude-<suffix>` entry from Keychain Access to remove the stored tokens; on Linux and Windows they are inside the `claude-subscriptions` directory removed above.

## Legacy zsh version

Earlier releases (up to `v0.3.0`) shipped a macOS-only Zsh implementation installed with `install.zsh`. It remains available at that tag but is no longer the recommended path; the cross-platform binary above supersedes it and is where new work happens.

## Limitations

- Usage values are cached and only update after a normal API response.
- Long-lived setup tokens cannot establish Claude Remote Control sessions.
- Claude Code does not currently provide a supported multi-profile `/login` interface.
- Inference-only setup tokens do not expose their account email through Claude's OAuth profile endpoint, so account identity must be established during browser authorization. Use **Edit** to set a display name (an email, for example) as the clearest identifier.
- On Windows, per-account config isolation falls back to a directory junction, hardlink, or copy when symbolic links are unavailable (i.e. without Developer Mode or admin rights).

## License

[MIT](LICENSE)
