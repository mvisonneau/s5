package cipher

import (
	"testing"

	"github.com/mvisonneau/s5/pkg/cipher/aes"
	"github.com/mvisonneau/s5/pkg/cipher/aws"
	"github.com/stretchr/testify/assert"
)

func TestNewAESClient(t *testing.T) {
	c, err := NewAESClient("cc6af4c2bf251c1cce0aebdbd39dc91d")
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.IsType(t, &aes.Client{}, c)

	c, err = NewAESClient("foo")
	assert.Error(t, err)
	assert.Nil(t, c)
}

func TestNewAWSClient(t *testing.T) {
	c, err := NewAWSClient("arn::kms::foo")
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.IsType(t, &aws.Client{}, c)
}

// func TestNewGCPClient(t *testing.T) {
// 	c, err := NewGCPClient("foo")
// 	assert.NoError(t, err)
// 	assert.NotNil(t, c)
// 	assert.IsType(t, &gcp.Client{}, c)
// }

func TestNewPGPClient(t *testing.T) {
	c, err := NewPGPClient("foo", "bar")
	assert.Error(t, err)
	assert.Nil(t, c)
	// TODO: Test with actual keys
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
	assert.NoError(t, err)
	assert.Equal(t, "abc", v)

	v, err = ParseInput("{{ s5:abc }}")
	assert.NoError(t, err)
	assert.Equal(t, "abc", v)

	v, err = ParseInput("{s5:abc}")
	assert.Error(t, err)
	assert.Equal(t, "", v)
}
