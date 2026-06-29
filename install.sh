#!/bin/sh
# Install the claude-accounts cross-platform binary from GitHub Releases.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/claude-code-tools/claude-code-account-switcher/v1.0.0/install.sh | sh
#
# Honors:
#   CLAUDE_ACCOUNTS_VERSION   release tag to install        (default below)
#   CLAUDE_ACCOUNTS_BIN_DIR   directory to install into     (default: a writable dir on PATH)
#
# This installer never touches your shell startup files, and never reads or
# transmits any stored tokens. macOS keeps tokens in the Keychain; Linux keeps
# them in a 0600 file under ~/.config/claude-subscriptions.

set -eu

REPO="claude-code-tools/claude-code-account-switcher"
DEFAULT_VERSION="v1.0.0"
VERSION="${CLAUDE_ACCOUNTS_VERSION:-$DEFAULT_VERSION}"

die() {
  printf 'error: %s\n' "$1" >&2
  exit 1
}

detect_os() {
  case "$(uname -s)" in
    Darwin) echo darwin ;;
    Linux) echo linux ;;
    *) die "unsupported OS $(uname -s); on Windows download the binary from https://github.com/$REPO/releases or use 'go install'" ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64 | amd64) echo amd64 ;;
    arm64 | aarch64) echo arm64 ;;
    *) die "unsupported architecture $(uname -m)" ;;
  esac
}

# Pick an install directory: an explicit override, else the first writable
# directory already on PATH (preferring ~/.local/bin), else ~/.local/bin.
choose_bin_dir() {
  if [ -n "${CLAUDE_ACCOUNTS_BIN_DIR:-}" ]; then
    echo "$CLAUDE_ACCOUNTS_BIN_DIR"
    return
  fi
  preferred="$HOME/.local/bin"
  case ":$PATH:" in
    *":$preferred:"*)
      if [ -d "$preferred" ] && [ -w "$preferred" ]; then
        echo "$preferred"
        return
      fi
      ;;
  esac
  IFS=:
  for dir in $PATH; do
    [ -n "$dir" ] || continue
    case "$dir" in
      */sbin | /usr/* | /bin | /opt/homebrew/* | /System/*) continue ;;
    esac
    if [ -d "$dir" ] && [ -w "$dir" ]; then
      unset IFS
      echo "$dir"
      return
    fi
  done
  unset IFS
  echo "$preferred"
}

main() {
  os="$(detect_os)"
  arch="$(detect_arch)"
  name="claude-accounts_${VERSION}_${os}_${arch}"
  url="https://github.com/$REPO/releases/download/$VERSION/${name}.tar.gz"

  command -v curl >/dev/null 2>&1 || die "curl is required"
  command -v tar >/dev/null 2>&1 || die "tar is required"

  tmp="$(mktemp -d)"
  trap 'rm -rf "$tmp"' EXIT

  printf 'Downloading claude-accounts %s (%s/%s)...\n' "$VERSION" "$os" "$arch"
  curl -fsSL "$url" -o "$tmp/archive.tar.gz" \
    || die "download failed: $url"

  # Verify the checksum when the release ships checksums.txt and a hasher exists.
  if curl -fsSL "https://github.com/$REPO/releases/download/$VERSION/checksums.txt" -o "$tmp/checksums.txt" 2>/dev/null; then
    expected="$(grep "${name}.tar.gz" "$tmp/checksums.txt" | awk '{print $1}')"
    if [ -n "$expected" ]; then
      if command -v shasum >/dev/null 2>&1; then
        actual="$(shasum -a 256 "$tmp/archive.tar.gz" | awk '{print $1}')"
      elif command -v sha256sum >/dev/null 2>&1; then
        actual="$(sha256sum "$tmp/archive.tar.gz" | awk '{print $1}')"
      fi
      [ -z "${actual:-}" ] || [ "$actual" = "$expected" ] \
        || die "checksum mismatch for ${name}.tar.gz"
    fi
  fi

  tar -C "$tmp" -xzf "$tmp/archive.tar.gz" claude-accounts \
    || die "could not extract claude-accounts from the archive"

  bin_dir="$(choose_bin_dir)"
  mkdir -p "$bin_dir" || die "cannot create install directory: $bin_dir"
  install -m 0755 "$tmp/claude-accounts" "$bin_dir/claude-accounts" \
    || die "cannot write to $bin_dir"

  printf '\nInstalled claude-accounts to %s\n' "$bin_dir/claude-accounts"
  case ":$PATH:" in
    *":$bin_dir:"*) ;;
    *) printf 'Note: %s is not on your PATH. Add it, then restart your shell.\n' "$bin_dir" ;;
  esac
  printf 'Run "claude-accounts" to add an account; each account also gets a claude-<suffix> launcher.\n'
}

main "$@"
