package cipher

import (
	cipherAWS "github.com/mvisonneau/s5/cipher/aws"

	log "github.com/sirupsen/logrus"
)

type AWS struct {
	client *cipherAWS.Client
}

func NewAWS(kmsKeyArn string) (*AWS, error) {
	c, err := cipherAWS.Init(
		&cipherAWS.Config{
			KmsKeyArn: kmsKeyArn,
		},
	)

	if err != nil {
		return nil, err
	}

	return &AWS{c}, nil
}

func (a *AWS) Cipher(value string) (string, error) {
	log.Debug("Ciphering using AWS KMS key")
	ciphered, err := a.client.Cipher(value)
	if err != nil {
		return "", err
	}
	return ciphered, nil
}

func (a *AWS) Decipher(value string) (string, error) {
	log.Debugf("Deciphering '%s' using AWS KMS key", value)
	plain, err := a.client.Decipher(value)
	if err != nil {
		return "", err
	}
	return plain, nil
}
