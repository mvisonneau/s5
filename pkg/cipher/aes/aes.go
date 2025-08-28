package aes

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"

	"github.com/mvisonneau/s5/internal/logs"
)

const (
	nonceByteLength = 12
)

// Config handles necessary information for AES.
type Config struct {
	// Key in a string format, usually passed from a CLI flag
	Key string
}

// Client is an handler for encryption functions.
type Client struct {
	cipher.AEAD
}

// NewClient configures a client for encryption purposes.
func NewClient(config *Config) (*Client, error) {
	key, err := hex.DecodeString(config.Key)
	if err != nil {
		return nil, errors.Wrap(err, "decoding hex value")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "initializing aes engine")
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "initializing gcm cipher engine")
	}

	return &Client{aead}, nil
}

// Cipher a value using the provided key.
func (c *Client) Cipher(ctx context.Context, value string) (string, error) {
	logs.LoggerFromContext(ctx).Debug("ciphering using AES")

	plaintext := []byte(value)

	nonce := make([]byte, nonceByteLength)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.Wrap(err, "generating nonce")
	}

	ciphertext := c.Seal(nil, nonce, plaintext, nil)

	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%x:%x", ciphertext, nonce))), nil
}

// Decipher a value using the TransitKey.
func (c *Client) Decipher(ctx context.Context, value string) (string, error) {
	logs.LoggerFromContext(ctx).Debug("deciphering using AES")

	str, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("base64decode - input: '%s'", value))
	}

	splittedStr := strings.Split(string(str), ":")

	ciphertext, err := hex.DecodeString(splittedStr[0])
	if err != nil {
		return "", errors.Wrap(err, "decoding ciphered string")
	}

	nonce, err := hex.DecodeString(splittedStr[1])
	if err != nil {
		return "", errors.Wrap(err, "decoding nonce from ciphered string")
	}

	plaintext, err := c.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.Wrap(err, "deciphering string with AES")
	}

	return string(plaintext), nil
}
