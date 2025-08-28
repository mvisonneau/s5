package gcp

import (
	"context"
	"encoding/base64"
	"fmt"

	cloudkms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/pkg/errors"

	"github.com/mvisonneau/s5/internal/logs"
)

// Config handles necessary information for AES.
type Config struct {
	KmsKeyName string
}

// Client is an handler for encryption functions.
type Client struct {
	*cloudkms.KeyManagementClient

	Context *context.Context
	Config  *Config
}

// NewClient configures a client for encryption purposes.
func NewClient(config *Config) (*Client, error) {
	ctx := context.Background()

	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "creating new client")
	}

	return &Client{
		client,
		&ctx,
		config,
	}, nil
}

// Cipher : Cipher a value using the provided key.
func (c *Client) Cipher(ctx context.Context, value string) (string, error) {
	logs.LoggerFromContext(ctx).Debug("ciphering using GCP KMS key")

	req := &kmspb.EncryptRequest{
		Name:      c.Config.KmsKeyName,
		Plaintext: []byte(value),
	}

	resp, err := c.Encrypt(*c.Context, req)
	if err != nil {
		return "", errors.Wrap(err, "ciphering value")
	}

	return base64.StdEncoding.EncodeToString(resp.GetCiphertext()), nil
}

// Decipher : Decipher a value using the TransitKey.
func (c *Client) Decipher(ctx context.Context, value string) (string, error) {
	logs.LoggerFromContext(ctx).Debug("deciphering using GCP KMS key")

	ciphertext, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("base64decode error, input : '%s'", value))
	}

	req := &kmspb.DecryptRequest{
		Name:       c.Config.KmsKeyName,
		Ciphertext: ciphertext,
	}

	resp, err := c.Decrypt(*c.Context, req)
	if err != nil {
		return "", errors.Wrap(err, "deciphering value")
	}

	return string(resp.GetPlaintext()), nil
}
