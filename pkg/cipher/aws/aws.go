package aws

import (
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Config handles necessary information for AES.
type Config struct {
	// Key in a string format, usually passed from a CLI flag
	KmsKeyArn string
}

// Client is an handler for encryption functions.
type Client struct {
	*kms.KMS
	Config *Config
}

// NewClient configures a client for encryption purposes.
func NewClient(config *Config) (*Client, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, errors.Wrap(err, "creating new AWS session")
	}

	return &Client{
		kms.New(sess),
		config,
	}, nil
}

// Cipher a value using the provided key.
func (c *Client) Cipher(value string) (string, error) {
	log.Debug("Ciphering using AWS KMS key")

	result, err := c.Encrypt(&kms.EncryptInput{
		KeyId:     aws.String(c.Config.KmsKeyArn),
		Plaintext: []byte(value),
	})
	if err != nil {
		return "", errors.Wrap(err, "ciphering using AWS KMS key")
	}

	return base64.StdEncoding.EncodeToString(result.CiphertextBlob), nil
}

// Decipher a value using the provided KMS Key.
func (c *Client) Decipher(value string) (string, error) {
	log.Debugf("Deciphering '%s' using AWS KMS key", value)

	ciphertext, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("base64decode - input: '%s'", value))
	}

	result, err := c.Decrypt(&kms.DecryptInput{CiphertextBlob: ciphertext})
	if err != nil {
		return "", errors.Wrap(err, "deciphering using AWS KMS key")
	}

	return string(result.Plaintext), nil
}
