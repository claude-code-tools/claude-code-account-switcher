// Package credstore reads a subscription's long-lived OAuth token from the OS
// credential store. The implementation is platform-specific: macOS uses the
// Keychain (via the security CLI, for parity with the zsh tool); Linux and
// Windows implementations are added per-OS.
//
// Get returns the token for a service name, or ("", nil) when no token is
// stored. A non-nil error indicates the store could not be queried at all.
package credstore
