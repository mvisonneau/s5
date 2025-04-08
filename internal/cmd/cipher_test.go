package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestCipher(t *testing.T) {
	// With a valid cipher engine
	ctx, flags, _ := NewTestContext()
	ctx.Command = &cli.Command{
		Name: "aes",
	}

	flags.String("key", "cc6af4c2bf251c1cce0aebdbd39dc91d", "")
	_ = flags.Parse([]string{"foo"})

	exitCode, err := Cipher(ctx)
	require.NoError(t, err)
	assert.Equal(t, 0, exitCode)

	// With a invalid cipher engine
	ctx.Command = &cli.Command{
		Name: "foo",
	}
	exitCode, err = Cipher(ctx)
	require.Error(t, err)
	assert.Equal(t, 1, exitCode)
}
