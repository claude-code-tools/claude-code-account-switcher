//go:build darwin

package credstore

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// Get reads the token stored in the login Keychain under the given service
// name. A missing item (security exit code 44, errSecItemNotFound) is reported
// as an empty token, not an error.
func Get(service string) (string, error) {
	cmd := exec.Command("/usr/bin/security", "find-generic-password",
		"-a", os.Getenv("USER"), "-s", service, "-w")
	out, err := cmd.Output()
	if err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			// 44 == errSecItemNotFound: no token configured for this service.
			return "", nil
		}
		return "", err
	}
	return strings.TrimRight(string(out), "\r\n"), nil
}
