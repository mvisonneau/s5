package transit

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/go-homedir"
)

// Config : Handles a Vault configuration
type Config struct {
	Key string
}

// Client is a Vault API client pointer
type Client struct {
	*api.Client
	Config *Config
}

// Init : Configure a Vault client and set a TransitKey to use
func Init(config *Config) (*Client, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return &Client{
		client,
		config,
	}, nil
}

// Cipher : Cipher a value using the TransitKey
func (c *Client) Cipher(value string) (string, error) {
	payload := make(map[string]interface{})
	payload["plaintext"] = base64.StdEncoding.EncodeToString([]byte(value))

	d, err := c.Logical().Write("transit/encrypt/"+c.Config.Key, payload)
	if err != nil {
		return "", fmt.Errorf("Vault error : %s", err)
	}

	re := regexp.MustCompile("(^vault:v1)")
	return re.ReplaceAllString(d.Data["ciphertext"].(string), "s5"), nil
}

// Decipher : Decipher a value using the TransitKey
func (c *Client) Decipher(value string) (string, error) {
	payload := make(map[string]interface{})
	re := regexp.MustCompile("(^s5:)")
	payload["ciphertext"] = re.ReplaceAllString(value, "vault:v1:")

	d, err := c.Logical().Write("transit/decrypt/"+c.Config.Key, payload)
	if err != nil {
		return "", fmt.Errorf("Vault error : %s", err)
	}

	output, err := base64.StdEncoding.DecodeString(d.Data["plaintext"].(string))
	return string(output), err
}

// getClient : Get a Vault client using Vault official params
func getClient() (*api.Client, error) {
	c, err := api.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating Vault client: %s", err.Error())
	}

	if len(os.Getenv("VAULT_ADDR")) == 0 {
		return nil, fmt.Errorf("VAULT_ADDR env is not defined")
	}

	c.SetAddress(os.Getenv("VAULT_ADDR"))

	token := os.Getenv("VAULT_TOKEN")
	if len(token) == 0 {
		home, _ := homedir.Dir()
		f, err := ioutil.ReadFile(home + "/.vault-token")
		if err != nil {
			return nil, fmt.Errorf("Vault token is not defined (VAULT_TOKEN or ~/.vault-token)")
		}

		token = string(f)
	}

	c.SetToken(token)

	return c, nil
}
