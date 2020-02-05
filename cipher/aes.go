package cipher

import (
	cipherAES "github.com/mvisonneau/s5/cipher/aes"

	log "github.com/sirupsen/logrus"
)

type AES struct {
	client *cipherAES.Client
}

func NewAES(key string) (*AES, error) {
	c, err := cipherAES.Init(
		&cipherAES.Config{
			Key: key,
		},
	)

	if err != nil {
		return nil, err
	}

	return &AES{c}, nil
}

func (a *AES) Cipher(value string) (string, error) {
	log.Debug("Ciphering using AES")
	ciphered, err := a.client.Cipher(value)
	if err != nil {
		return "", err
	}
	return ciphered, nil
}

func (a *AES) Decipher(value string) (string, error) {
	log.Debugf("Deciphering '%s' using AES", value)
	plain, err := a.client.Decipher(value)
	if err != nil {
		return "", err
	}
	return plain, nil
}
