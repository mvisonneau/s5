package cipher

import (
	cipherVault "github.com/mvisonneau/s5/cipher/vault"

	log "github.com/sirupsen/logrus"
)

type Vault struct {
	client *cipherVault.Client
}

func NewVault(transitKey string) (*Vault, error) {
	c, err := cipherVault.Init(
		&cipherVault.Config{
			Key: transitKey,
		},
	)

	if err != nil {
		return nil, err
	}

	return &Vault{c}, nil
}

func (v *Vault) cipher(value string) (string, error) {
	log.Debug("Ciphering using Vault transit key")
	ciphered, err := v.client.Cipher(value)
	if err != nil {
		return "", err
	}
	return ciphered, nil
}

func (v *Vault) Decipher(value string) (string, error) {
	log.Debugf("Deciphering '%s' using Vault transit key", value)
	plain, err := v.client.Decipher(value)
	if err != nil {
		return "", err
	}
	return plain, nil
}
