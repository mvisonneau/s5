package aws

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/smithy-go/logging"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/mvisonneau/s5/internal/app"
	"github.com/mvisonneau/s5/internal/logs"
)

// Config handles necessary information for AES.
type Config struct {
	// Key in a string format, usually passed from a CLI flag
	KmsKeyArn string
}

// Client is a handler for encryption functions.
type Client struct {
	*kms.Client

	Config *Config
}

// NewClient configures a client for encryption purposes.
func NewClient(ctx context.Context, cfg *Config) (*Client, error) {
	logger := logging.LoggerFunc(func(_ logging.Classification, _ string, entries ...interface{}) {
		for _, entry := range entries {
			logs.LoggerFromContext(ctx).Debug("", zap.Any("event", entry))
		}
	})

	awsConfig, err := config.LoadDefaultConfig(
		ctx,
		config.WithLogger(logger),
		config.WithAppID(app.Name),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		kms.NewFromConfig(awsConfig),
		cfg,
	}, nil
}

// Cipher a value using the provided key.
func (c *Client) Cipher(ctx context.Context, value string) (string, error) {
	logs.LoggerFromContext(ctx).Debug("ciphering using AWS KMS key")

	result, err := c.Encrypt(ctx, &kms.EncryptInput{
		KeyId:     aws.String(c.Config.KmsKeyArn),
		Plaintext: []byte(value),
	})
	if err != nil {
		return "", errors.Wrap(err, "ciphering using AWS KMS key")
	}

	return base64.StdEncoding.EncodeToString(result.CiphertextBlob), nil
}

// Decipher a value using the provided KMS Key.
func (c *Client) Decipher(ctx context.Context, value string) (string, error) {
	logs.LoggerFromContext(ctx).Debug("deciphering using AWS KMS key")

	ciphertext, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("base64decode - input: '%s'", value))
	}

	result, err := c.Decrypt(ctx, &kms.DecryptInput{CiphertextBlob: ciphertext})
	if err != nil {
		return "", errors.Wrap(err, "deciphering using AWS KMS key")
	}

	return string(result.Plaintext), nil
}
