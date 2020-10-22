package aes

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)
	assert.NotNil(t, c)

	c, err = NewClient(invalidTestConfig)
	assert.Error(t, err)
	assert.Nil(t, c)
}

func TestValidCipherDecipher(t *testing.T) {
	c, err := NewClient(validTestConfig)
	assert.NoError(t, err)

	ciphered, err := c.Cipher(valueToCipher)
	assert.NoError(t, err)
	assert.NotEqual(t, valueToCipher, ciphered)

	deciphered, err := c.Decipher(ciphered)
	assert.NoError(t, err)
	assert.Equal(t, valueToCipher, deciphered)
}

func TestInvalidDecipher(t *testing.T) {
	c, err := NewClient(validTestConfig)
	assert.NoError(t, err)

	deciphered, err := c.Decipher("not_a_valid_ciphered_value")
	assert.Error(t, err)
	assert.Equal(t, "", deciphered)
}
