package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

// Config handles necessary information for AES
type Config struct {
	// Key in a string format, usually passed from a CLI flag
	Key string
}

// Client is an handler for encryption functions
type Client struct {
	cipher.AEAD
}

// Init : Configures a client for encryption purposes
func Init(config *Config) (*Client, error) {
	key, err := hex.DecodeString(config.Key)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &Client{aead}, nil
}

// Cipher : Cipher a value using the provided key
func (c *Client) Cipher(value string) (string, error) {
	plaintext := []byte(value)

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := c.Seal(nil, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%x:%x", ciphertext, nonce))), nil
}

// Decipher : Decipher a value using the TransitKey
func (c *Client) Decipher(value string) (string, error) {
	str, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", fmt.Errorf("base64decode error : %s - value : %s", err, value)
	}

	s := strings.Split(string(str), ":")

	ciphertext, err := hex.DecodeString(s[0])
	if err != nil {
		return "", err
	}

	nonce, err := hex.DecodeString(s[1])
	if err != nil {
		return "", err
	}

	plaintext, err := c.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
