package cipher

import (
	cipherGCP "github.com/mvisonneau/s5/cipher/gcp"

	log "github.com/sirupsen/logrus"
)

type GCP struct {
	client *cipherGCP.Client
}

func NewGCP(kmsKeyName string) (*GCP, error) {
	c, err := cipherGCP.Init(
		&cipherGCP.Config{
			KmsKeyName: kmsKeyName,
		},
	)

	if err != nil {
		return nil, err
	}

	return &GCP{c}, nil
}

func (g *GCP) Cipher(value string) (string, error) {
	log.Debug("Ciphering using a GCP KMS key")
	ciphered, err := g.client.Cipher(value)
	if err != nil {
		return "", err
	}
	return ciphered, nil
}

func (g *GCP) Decipher(value string) (string, error) {
	log.Debugf("Deciphering '%s' using a GCP KMS key", value)
	plain, err := g.client.Decipher(value)
	if err != nil {
		return "", err
	}
	return plain, nil
}
