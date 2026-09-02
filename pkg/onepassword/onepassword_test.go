package onepassword

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	helperProcessEnv = "S5_WANT_HELPER_PROCESS"
	validReference   = "op://Personal/s5/aes-key"
)

// TestHelperProcess is not a real test, it impersonates the `op` CLI when
// executed through fakeCLI().
func TestHelperProcess(t *testing.T) {
	if os.Getenv(helperProcessEnv) != "1" {
		return
	}

	_, _ = fmt.Fprint(os.Stdout, os.Getenv("S5_HELPER_STDOUT"))
	_, _ = fmt.Fprint(os.Stderr, os.Getenv("S5_HELPER_STDERR"))

	exitCode, _ := strconv.Atoi(os.Getenv("S5_HELPER_EXIT_CODE"))

	os.Exit(exitCode)
}

// fakeCLI replaces the `op` CLI with the test binary itself, which returns the
// provided stdout, stderr and exit code. It also records the arguments the CLI
// got called with.
func fakeCLI(t *testing.T, stdout, stderr string, exitCode int) *[]string {
	t.Helper()

	var args []string

	original := execCommandContext

	execCommandContext = func(ctx context.Context, name string, arg ...string) *exec.Cmd {
		args = append([]string{name}, arg...)

		// #nosec G204 -- we are executing the test binary itself
		c := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestHelperProcess$", "-test.v=false")
		c.Env = append(os.Environ(),
			helperProcessEnv+"=1",
			"S5_HELPER_STDOUT="+stdout,
			"S5_HELPER_STDERR="+stderr,
			"S5_HELPER_EXIT_CODE="+strconv.Itoa(exitCode),
		)

		return c
	}

	t.Cleanup(func() {
		execCommandContext = original
	})

	return &args
}

func TestIsSecretReference(t *testing.T) {
	assert.True(t, IsSecretReference(validReference))
	assert.True(t, IsSecretReference("op://Personal/s5/section/aes-key"))
	assert.False(t, IsSecretReference("cc6af4c2bf251c1cce0aebdbd39dc91d"))
	assert.False(t, IsSecretReference(""))
	assert.False(t, IsSecretReference("/op://Personal/s5/aes-key"))
}

func TestReadInvalidReference(t *testing.T) {
	args := fakeCLI(t, "foo", "", 0)

	for _, reference := range []string{
		"op://",
		"op://Personal",
		"op://Personal/s5",
		"op://Personal/s5/section/subsection/aes-key",
		"op://Personal/s5/ ",
		"op://Personal//aes-key",
	} {
		value, err := Read(context.TODO(), reference)
		require.ErrorContains(t, err, "invalid 1password secret reference")
		assert.Empty(t, value)
	}

	// the CLI should not have been executed at all
	assert.Empty(t, *args)
}

func TestReadReferenceWithSpacesAndQueryParameters(t *testing.T) {
	args := fakeCLI(t, "cc6af4c2bf251c1cce0aebdbd39dc91d", "", 0)

	reference := "op://Private Vault/my s5 item/section/one-time password?attribute=otp"

	value, err := Read(context.TODO(), reference)
	require.NoError(t, err)
	assert.Equal(t, "cc6af4c2bf251c1cce0aebdbd39dc91d", value)
	assert.Equal(t, []string{CLI, "read", "--no-newline", reference}, *args)
}

func TestReadValidReference(t *testing.T) {
	args := fakeCLI(t, "cc6af4c2bf251c1cce0aebdbd39dc91d\n", "", 0)

	value, err := Read(context.TODO(), validReference)
	require.NoError(t, err)
	assert.Equal(t, "cc6af4c2bf251c1cce0aebdbd39dc91d", value)
	assert.Equal(t, []string{CLI, "read", "--no-newline", validReference}, *args)
}

func TestReadEmptyValue(t *testing.T) {
	fakeCLI(t, "  \n", "", 0)

	value, err := Read(context.TODO(), validReference)
	require.ErrorContains(t, err, "1password returned an empty value")
	assert.Empty(t, value)
}

func TestReadCLIError(t *testing.T) {
	fakeCLI(t, "", "[ERROR] could not read secret: not found\n", 1)

	value, err := Read(context.TODO(), validReference)
	require.ErrorContains(t, err, "could not read secret: not found")
	assert.Empty(t, value)
}

func TestReadCLIErrorWithoutStderr(t *testing.T) {
	fakeCLI(t, "", "", 1)

	value, err := Read(context.TODO(), validReference)
	require.ErrorContains(t, err, "reading '"+validReference+"' from 1password")
	assert.Empty(t, value)
}

func TestReadCLINotFound(t *testing.T) {
	original := execCommandContext

	execCommandContext = func(ctx context.Context, _ string, arg ...string) *exec.Cmd {
		return exec.CommandContext(ctx, "s5-tests-op-cli-which-does-not-exist", arg...)
	}

	t.Cleanup(func() {
		execCommandContext = original
	})

	value, err := Read(context.TODO(), validReference)
	require.ErrorContains(t, err, "the 'op' cli is required")
	assert.Empty(t, value)
}
