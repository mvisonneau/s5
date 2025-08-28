package vault

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"
	"regexp"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"

	"github.com/mvisonneau/s5/internal/logs"
)

// Config : Handles a Vault configuration.
type Config struct {
	Key string
}

// Client is a Vault API client pointer.
type Client struct {
	*api.Client

	Config *Config
}

// NewClient configures a Vault client and set a TransitKey to use.
func NewClient(config *Config) (*Client, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return &Client{
		client,
		config,
	}, nil
}

// Cipher : Cipher a value using the TransitKey.
func (c *Client) Cipher(ctx context.Context, value string) (string, error) {
	logs.LoggerFromContext(ctx).Debug("ciphering using a Vault transit key")

	payload := make(map[string]interface{})
	payload["plaintext"] = base64.StdEncoding.EncodeToString([]byte(value))

	d, err := c.Logical().Write("transit/encrypt/"+c.Config.Key, payload)
	if err != nil {
		return "", errors.Wrap(err, "vault client")
	}

	re := regexp.MustCompile("(^vault:v1:)")

	return re.ReplaceAllString(d.Data["ciphertext"].(string), ""), nil
}

// Decipher : Decipher a value using the TransitKey.
func (c *Client) Decipher(ctx context.Context, value string) (string, error) {
	logs.LoggerFromContext(ctx).Debug("deciphering using a Vault transit key")

	payload := make(map[string]interface{})
	payload["ciphertext"] = "vault:v1:" + value

	d, err := c.Logical().Write("transit/decrypt/"+c.Config.Key, payload)
	if err != nil {
		return "", errors.Wrap(err, "vault client")
	}

	output, err := base64.StdEncoding.DecodeString(d.Data["plaintext"].(string))

	return string(output), err
}

// getClient : Get a Vault client using Vault official params.
func getClient() (*api.Client, error) {
	c, err := api.NewClient(nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating vault client")
	}

	if len(os.Getenv("VAULT_ADDR")) == 0 {
		return nil, errors.New("VAULT_ADDR env is not defined")
	}

	if err = c.SetAddress(os.Getenv("VAULT_ADDR")); err != nil {
		return nil, errors.Wrap(err, "configuring vault endpoint address")
	}

	token := os.Getenv("VAULT_TOKEN")
	if len(token) == 0 {
		home, _ := homedir.Dir()

		f, err := os.ReadFile(filepath.Clean(home + "/.vault-token"))
		if err != nil {
			return nil, errors.New("vault token is not defined (VAULT_TOKEN or ~/.vault-token)")
		}

		token = string(f)
	}

	c.SetToken(token)

	return c, nil
}
