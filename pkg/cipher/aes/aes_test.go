package aes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	valueToCipher = "foo"
)

var validTestConfig = &Config{
	Key: "cc6af4c2bf251c1cce0aebdbd39dc91d",
}

var invalidTestConfig = &Config{
	Key: "bar",
}

func TestNewClient(t *testing.T) {
	c, err := NewClient(validTestConfig)
	require.NoError(t, err)
	assert.NotNil(t, c)

	c, err = NewClient(invalidTestConfig)
	require.Error(t, err)
	assert.Nil(t, c)
}

func TestValidCipherDecipher(t *testing.T) {
	c, err := NewClient(validTestConfig)
	require.NoError(t, err)

	ciphered, err := c.Cipher(context.TODO(), valueToCipher)
	require.NoError(t, err)
	assert.NotEqual(t, valueToCipher, ciphered)

	deciphered, err := c.Decipher(context.TODO(), ciphered)
	require.NoError(t, err)
	assert.Equal(t, valueToCipher, deciphered)
}

func TestInvalidDecipher(t *testing.T) {
	c, err := NewClient(validTestConfig)
	require.NoError(t, err)

	deciphered, err := c.Decipher(context.TODO(), "not_a_valid_ciphered_value")
	require.Error(t, err)
	assert.Empty(t, deciphered)
}
