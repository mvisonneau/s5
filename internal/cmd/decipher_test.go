package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestDecipher(t *testing.T) {
	// Valid cipher engine with valid input
	ctx, flags, _ := NewTestContext()
	ctx.Command = &cli.Command{
		Name: "aes",
	}
	flags.String("key", "cc6af4c2bf251c1cce0aebdbd39dc91d", "")
	flags.Parse([]string{"{{s5:MmZmZTI0NDI1NjY3YTdhNjZhZjFmMGZjMzdkZjM0OTBiZGY0MDc6YTEzNzdlOGJkMTc2ZDg5NjE2ZTJlNjll}}"})

	exitCode, err := Decipher(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 0, exitCode)

	// Valid cipher engine with invalid input
	flags.Parse([]string{"{{s5:bar}}"})
	exitCode, err = Decipher(ctx)
	assert.Error(t, err)
	assert.Equal(t, 1, exitCode)

	// With a invalid cipher engine
	ctx.Command = &cli.Command{
		Name: "foo",
	}
	exitCode, err = Decipher(ctx)
	assert.Error(t, err)
	assert.Equal(t, 1, exitCode)
}
