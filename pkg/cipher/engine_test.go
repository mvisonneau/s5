package cipher

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mvisonneau/s5/pkg/cipher/aes"
	"github.com/mvisonneau/s5/pkg/cipher/aws"
)

func TestNewAESClient(t *testing.T) {
	c, err := NewAESClient("cc6af4c2bf251c1cce0aebdbd39dc91d")
	require.NoError(t, err)
	assert.NotNil(t, c)
	assert.IsType(t, &aes.Client{}, c)

	c, err = NewAESClient("foo")
	require.Error(t, err)
	assert.Nil(t, c)
}

func TestNewAWSClient(t *testing.T) {
	c, err := NewAWSClient(context.TODO(), "arn::kms::foo")
	require.NoError(t, err)
	assert.NotNil(t, c)
	assert.IsType(t, &aws.Client{}, c)
}

// func TestNewGCPClient(t *testing.T) {
// 	c, err := NewGCPClient("foo")
// 	assert.NoError(t, err)
// 	assert.NotNil(t, c)
// 	assert.IsType(t, &gcp.Client{}, c)
// }

// TODO: Test with actual keys.
func TestNewPGPClient(t *testing.T) {
	c, err := NewPGPClient("foo", "bar")
	require.Error(t, err)
	assert.Nil(t, c)
}

// func TestNewVaultClient(t *testing.T) {
// 	c, err := NewVaultClient("foo")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "foo", c.Config.Key)
// 	assert.IsType(t, &vault.Client{}, c)
// }

func TestGenerateOutput(t *testing.T) {
	assert.Equal(t, "{{s5:foo}}", GenerateOutput("foo"))
	assert.Equal(t, "{{s5:bar}}", GenerateOutput("bar"))
}

func TestParseInput(t *testing.T) {
	v, err := ParseInput("{{s5:abc}}")
	require.NoError(t, err)
	assert.Equal(t, "abc", v)

	v, err = ParseInput("{{ s5:abc }}")
	require.NoError(t, err)
	assert.Equal(t, "abc", v)

	v, err = ParseInput("{s5:abc}")
	require.Error(t, err)
	assert.Empty(t, v)
}
