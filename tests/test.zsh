#!/bin/zsh

emulate -L zsh
set -euo pipefail

readonly PROJECT_ROOT="${0:A:h:h}"
readonly TEST_ROOT="$(mktemp -d "${TMPDIR:-/tmp}/claude-accounts-test.XXXXXX")"
trap 'rm -rf "$TEST_ROOT"' EXIT INT TERM

# Standalone mode skips function generation; a leaked value would make the
# runtime tests below fail spuriously. Keep the suite hermetic.
unset CLAUDE_ACCOUNTS_STANDALONE CLAUDE_ACCOUNTS_RUNTIME

fail() {
  print -u2 -r -- "test failure: $*"
  exit 1
}

export HOME="$TEST_ROOT/home"
export XDG_CONFIG_HOME="$HOME/.config"
mkdir -p "$HOME" "$XDG_CONFIG_HOME/claude-subscriptions/usage"
mkdir -p "$HOME/bin"
export PATH="$HOME/bin:/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"

export CLAUDE_SUBSCRIPTIONS_FILE="$XDG_CONFIG_HOME/claude-subscriptions/accounts.tsv"
export CLAUDE_SUBSCRIPTIONS_USAGE_DIR="$XDG_CONFIG_HOME/claude-subscriptions/usage"
export CLAUDE_SUBSCRIPTIONS_USAGE_SETTINGS="$XDG_CONFIG_HOME/claude-subscriptions/usage-settings.json"
export CLAUDE_ACCOUNTS_BIN_DIR="$HOME/bin"

print -r -- $'gmail\tGmail Work\tClaude Code Subscription: gmail' > "$CLAUDE_SUBSCRIPTIONS_FILE"
print -r -- '#!/bin/sh' > "$CLAUDE_ACCOUNTS_BIN_DIR/claude-accounts"
chmod 700 "$CLAUDE_ACCOUNTS_BIN_DIR/claude-accounts"
print -r -- '{"captured_at":1782540000,"rate_limits":{"five_hour":{"used_percentage":23.5},"seven_day":{"used_percentage":41.2}}}' > "$CLAUDE_SUBSCRIPTIONS_USAGE_DIR/gmail.json"

source "$PROJECT_ROOT/src/claude-accounts.zsh"

typeset -f claude-accounts >/dev/null || fail "claude-accounts function was not loaded"
typeset -f claude-gmail >/dev/null || fail "claude-gmail function was not generated"
[[ -L "$CLAUDE_ACCOUNTS_BIN_DIR/claude-gmail" ]] || fail "claude-gmail executable was not generated"

if typeset -f claude-claude-gmail >/dev/null; then
  print -u2 -r -- "unexpected duplicate claude- prefix"
  exit 1
fi

_claude_subscription_suggest_slug "claude-naver"
[[ "$REPLY" == "naver" ]] || fail "claude- prefix was not normalized"

_claude_subscription_usage_summary "gmail"
[[ "$REPLY" == "5h 24% · 7d 41% used" ]] || fail "usage summary was not formatted correctly"

# capture-usage.zsh renders the pinned account label on the status line.
statusline_out="$(CLAUDE_SUBSCRIPTION_SLUG="statusliner" CLAUDE_SUBSCRIPTION_LABEL="Status Liner" \
  zsh "$PROJECT_ROOT/src/capture-usage.zsh" <<< '{"rate_limits":{"five_hour":{"used_percentage":1}}}')"
[[ "$statusline_out" == "Status Liner" ]] || \
  fail "status line did not render the account label (got: ${statusline_out})"

# It falls back to the slug when no label is exported (sessions launched before upgrade).
statusline_out="$(CLAUDE_SUBSCRIPTION_SLUG="statusliner" \
  zsh "$PROJECT_ROOT/src/capture-usage.zsh" <<< '{}')"
[[ "$statusline_out" == "statusliner" ]] || \
  fail "status line did not fall back to the slug (got: ${statusline_out})"

# It stays silent for sessions not launched through the switcher.
statusline_out="$(CLAUDE_SUBSCRIPTION_SLUG="" \
  zsh "$PROJECT_ROOT/src/capture-usage.zsh" <<< '{}')"
[[ -z "$statusline_out" ]] || \
  fail "status line should be empty without a subscription slug (got: ${statusline_out})"

if _claude_subscription_has_explicit_session_name --continue; then
  fail "unnamed sessions should receive the account name"
fi
_claude_subscription_has_explicit_session_name --name "review" || \
  fail "--name was not detected"
_claude_subscription_has_explicit_session_name --name=review || \
  fail "--name=value was not detected"
_claude_subscription_has_explicit_session_name -n "review" || \
  fail "-n was not detected"
if _claude_subscription_has_explicit_session_name -- --name "prompt text"; then
  fail "arguments after -- should not be parsed as session options"
fi

# Per-account config isolation strips the cached account but shares the rest.
mkdir -p "$HOME/.claude"
print -r -- '{"oauthAccount":{"organizationUuid":"gmail-org"},"numStartups":7,"projects":{}}' > "$HOME/.claude.json"
print -r -- '{"model":"opus"}' > "$HOME/.claude/settings.json"
print -r -- 'lock' > "$HOME/.claude/daemon.lock"
if _claude_subscription_config_dir gmail; then
  acct="$REPLY"
  [[ -f "$acct/.claude.json" && ! -L "$acct/.claude.json" ]] || fail "per-account .claude.json should be a real file"
  if command -v jq >/dev/null 2>&1; then
    [[ "$(jq -r '.oauthAccount // "gone"' "$acct/.claude.json")" == "gone" ]] || fail "oauthAccount was not stripped"
    [[ "$(jq -r '.numStartups' "$acct/.claude.json")" == "7" ]] || fail "shared keys were not preserved"
  fi
  [[ -L "$acct/settings.json" && "$(readlink "$acct/settings.json")" == "$HOME/.claude/settings.json" ]] || \
    fail "settings.json was not shared via symlink"
  [[ ! -e "$acct/daemon.lock" ]] || fail "daemon.lock should not be mirrored into the account config"
else
  print -r -- "note: account isolation unavailable (no jq/python3); skipping isolation assertions"
fi

print -r -- "runtime tests passed"
