#!/bin/zsh

set -u

slug="${CLAUDE_SUBSCRIPTION_SLUG:-}"
[[ -n "$slug" && "$slug" != -* && "$slug" != *- && "$slug" != *[^a-z0-9-]* ]] || exit 0

# Drain the status line payload so the producer never observes a broken pipe.
input="$(cat)"

# Cache usage when jq is available and the payload carries rate limits.
if (( ${+commands[jq]} )); then
  rate_limits="$(printf '%s' "$input" | jq -c '.rate_limits // empty' 2>/dev/null)"
  if [[ -n "$rate_limits" ]]; then
    usage_dir="${XDG_CONFIG_HOME:-$HOME/.config}/claude-subscriptions/usage"
    if mkdir -p "$usage_dir"; then
      chmod 700 "$usage_dir" 2>/dev/null

      usage_file="$usage_dir/${slug}.json"
      tmp_file="${usage_file}.tmp.$$"
      captured_at="$(date +%s)"

      if printf '%s' "$input" | jq -c \
        --argjson captured_at "$captured_at" \
        '{captured_at: $captured_at, rate_limits: .rate_limits}' > "$tmp_file" 2>/dev/null; then
        chmod 600 "$tmp_file" 2>/dev/null
        mv "$tmp_file" "$usage_file"
      else
        rm -f "$tmp_file"
      fi
    fi
  fi
fi

# Always render the pinned account so the active subscription stays visible for
# the whole session, regardless of any in-session /login. The label is pinned at
# launch and reflects the injected token, not Claude's mutable /login state.
label="${CLAUDE_SUBSCRIPTION_LABEL:-$slug}"
label="${label//[$'\t\n\r']/ }"
print -r -- "$label"
