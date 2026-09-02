// Package onepassword resolves 1Password secret references, allowing secrets
// used by s5 to be kept in 1Password rather than in environment variables or
// files. References are resolved by shelling out to the `op` CLI, which means
// they can be unlocked interactively (biometrics, desktop app integration) or
// non-interactively through a service account token.
package onepassword

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/term"

	"github.com/mvisonneau/s5/internal/logs"
)

const (
	// SecretReferenceScheme is the prefix of a 1Password secret reference.
	SecretReferenceScheme string = "op://"

	// CLI is the name of the 1Password command line binary.
	CLI string = "op"

	// segment is a single, non blank, component of a secret reference. Vault,
	// item and field names may contain spaces and fields can hold query
	// parameters, eg: op://prod/db/one-time password?attribute=otp.
	segment string = `[^/\r\n]*[^\s/][^/\r\n]*`
)

// SecretReferenceRegexp is defining the syntax of a 1Password secret
// reference: op://<vault>/<item>[/<section>]/<field>.
var SecretReferenceRegexp = regexp.MustCompile(`^op://` + segment + `(/` + segment + `){2,3}$`)

// execCommandContext is a variable in order to be able to fake the `op` CLI
// within the tests.
var execCommandContext = exec.CommandContext

// IsSecretReference returns whether a value is a 1Password secret reference.
func IsSecretReference(value string) bool {
	return strings.HasPrefix(value, SecretReferenceScheme)
}

// Read resolves a 1Password secret reference and returns its value.
func Read(ctx context.Context, reference string) (string, error) {
	logs.LoggerFromContext(ctx).Debug(
		"resolving a 1password secret reference",
		zap.String("reference", reference),
	)

	if !SecretReferenceRegexp.MatchString(reference) {
		return "", errors.Errorf(
			"invalid 1password secret reference '%s', should be '%s<vault>/<item>[/<section>]/<field>'",
			reference,
			SecretReferenceScheme,
		)
	}

	var stdout, stderr bytes.Buffer

	// #nosec G204 -- the reference is validated above and passed as an argument,
	// no shell is involved.
	c := execCommandContext(ctx, CLI, "read", "--no-newline", reference)
	c.Stdout = &stdout

	if term.IsTerminal(int(os.Stdin.Fd())) {
		// hand over the terminal so that the CLI remains able to prompt for
		// credentials.
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
	} else {
		// the CLI must not consume the stdin of s5, which may hold the content
		// to (de)cipher.
		c.Stdin = nil
		c.Stderr = &stderr
	}

	if err := c.Run(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return "", errors.Wrapf(err, "the '%s' cli is required in order to resolve 1password secret references", CLI)
		}

		if msg := strings.TrimSpace(stderr.String()); len(msg) > 0 {
			return "", errors.Wrapf(err, "reading '%s' from 1password: %s", reference, msg)
		}

		return "", errors.Wrapf(err, "reading '%s' from 1password", reference)
	}

	value := strings.TrimSpace(stdout.String())
	if len(value) == 0 {
		return "", errors.Errorf("1password returned an empty value for '%s'", reference)
	}

	return value, nil
}
