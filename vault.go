package main

import (
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/vault/api"
)

// Vault : Handles a Vault API Client and a TransitKey name
type Vault struct {
	Client *api.Client
	TransitKey string
}

// Configure : Configure a Vault client and set a TransitKey to use
func (v *Vault) Configure(address, token, transitKey string) error {
	v.SetTransitKey(transitKey)
	return v.ConfigureClient(address, token)
}

// ConfigureClient : Configure a Vault client and set a TransitKey to use
func (v *Vault) ConfigureClient(address, token string) error {
	var err error
	v.Client, err = api.NewClient(nil)
	if err != nil {
		return fmt.Errorf("Error creating Vault client: %s", err.Error())
	}

	if len(address) == 0 {
		return fmt.Errorf("Vault address is not defined")
	}

	if len(token) == 0 {
		return fmt.Errorf("Vault token is not defined")
	}

	v.Client.SetAddress(address)
	v.Client.SetToken(token)

	return nil
}

func (v *Vault) SetTransitKey(transitKey string) {
	v.TransitKey = transitKey
}

// Cipher : Cipher a value using the TransitKey
func (v *Vault) Cipher(value string) (string, error) {
	payload := make(map[string]interface{})
	payload["plaintext"] = base64.StdEncoding.EncodeToString([]byte(value))

	d, err := v.Client.Logical().Write("transit/encrypt/"+v.TransitKey, payload)
	if err != nil {
		return "", fmt.Errorf("Vault error : %s", err )
	}

	return d.Data["ciphertext"].(string), nil
}

// Decipher : Decipher a value using the TransitKey
func (v *Vault) Decipher(value string) (string, error) {
	payload := make(map[string]interface{})
	payload["ciphertext"] = value

	d, err := v.Client.Logical().Write("transit/decrypt/"+v.TransitKey, payload)
	if err != nil {
		return "", fmt.Errorf("Vault error : %s", err)
	}

	output, err := base64.StdEncoding.DecodeString(d.Data["plaintext"].(string))
	return string(output), err
}
