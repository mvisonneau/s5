package gcp

import (
	"context"
	"encoding/base64"
	"fmt"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"

	log "github.com/sirupsen/logrus"
)

// Config handles necessary information for AES
type Config struct {
	KmsKeyName string
}

// Client is an handler for encryption functions
type Client struct {
	*cloudkms.KeyManagementClient
	Context *context.Context
	Config  *Config
}

// NewClient configures a client for encryption purposes
func NewClient(config *Config) (*Client, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	return &Client{
		client,
		&ctx,
		config,
	}, nil
}

// Cipher : Cipher a value using the provided key
func (c *Client) Cipher(value string) (string, error) {
	log.Debug("Ciphering using a GCP KMS key")
	req := &kmspb.EncryptRequest{
		Name:      c.Config.KmsKeyName,
		Plaintext: []byte(value),
	}

	resp, err := c.Encrypt(*c.Context, req)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(resp.Ciphertext), nil
}

// Decipher : Decipher a value using the TransitKey
func (c *Client) Decipher(value string) (string, error) {
	log.Debugf("Deciphering '%s' using a GCP KMS key", value)
	ciphertext, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", fmt.Errorf("base64decode error : %s - value : %s", err, value)
	}

	req := &kmspb.DecryptRequest{
		Name:       c.Config.KmsKeyName,
		Ciphertext: ciphertext,
	}

	resp, err := c.Decrypt(*c.Context, req)
	if err != nil {
		return "", err
	}

	return string(resp.Plaintext), nil
}
