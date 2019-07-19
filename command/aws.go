package command

import (
	cipherAWS "github.com/mvisonneau/s5/cipher/aws"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type aws struct {
	client *cipherAWS.Client
}

func (a *aws) configure(ctx *cli.Context) (err error) {
	a.client, err = cipherAWS.Init(
		&cipherAWS.Config{
			KmsKeyArn: ctx.String("kms-key-arn"),
		},
	)
	return
}

func (a *aws) cipher(value string) (string, error) {
	log.Debug("Ciphering using AWS KMS key")
	ciphered, err := a.client.Cipher(value)
	if err != nil {
		return "", err
	}
	return ciphered, nil
}

func (a *aws) decipher(value string) (string, error) {
	log.Debugf("Deciphering '%s' using AWS KMS key", value)
	plain, err := a.client.Decipher(value)
	if err != nil {
		return "", err
	}
	return plain, nil
}
