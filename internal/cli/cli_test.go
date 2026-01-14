package cli

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func NewTestApp() (app *cli.Command) {
	app = NewApp()
	app.Writer = io.Discard
	app.ErrWriter = io.Discard

	return
}

func TestRun(t *testing.T) {
	assert.NotPanics(t, func() {
		_ = NewTestApp().Run(context.Background(), os.Args)
	})
}

func TestNewApp(t *testing.T) {
	app := NewTestApp()
	assert.Equal(t, "s5", app.Name)
	assert.Equal(t, "v0.0.0-dev", app.Version)
}

func TestCipher(t *testing.T) {
	app := NewTestApp()

	tests := []struct {
		name      string
		arguments []string
		expectErr string
	}{
		{
			name:      "valid",
			arguments: []string{"s5", "cipher", "aes", "--key", "cc6af4c2bf251c1cce0aebdbd39dc91d", "coucou"},
			expectErr: "",
		},
		{
			name:      "invalid",
			arguments: []string{"s5", "cipher", "aes", "--key"},
			expectErr: "flag needs an argument: --key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.arguments

			err := app.Run(context.Background(), os.Args)
			if tt.expectErr != "" {
				require.ErrorContains(t, err, tt.expectErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDecipher(t *testing.T) {
	app := NewTestApp()

	tests := []struct {
		name      string
		arguments []string
		expectErr string
	}{
		{
			name: "valid",
			arguments: []string{
				"s5", "decipher", "aes", "--key",
				"cc6af4c2bf251c1cce0aebdbd39dc91d",
				"{{s5:MTdkNjc2OWRlNDliMTljZjAwNGM3YjI3MTRkNGRjNWFmZGE0NDUzYWM5Zjg6M2NhYzZjNjliMTQxZDdmZGU5MmEyZTMy}}",
			},
			expectErr: "",
		},
		{
			name:      "invalid",
			arguments: []string{"s5", "decipher", "aes"},
			expectErr: "crypto/aes: invalid key size 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.arguments
			err := app.Run(context.Background(), os.Args)
			if tt.expectErr != "" {
				require.ErrorContains(t, err, tt.expectErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
