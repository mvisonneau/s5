package command

import (
	cipherAES "github.com/mvisonneau/s5/cipher/aes"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type aes struct {
	client *cipherAES.Client
}

func (a *aes) configure(ctx *cli.Context) (err error) {
	a.client, err = cipherAES.Init(
		&cipherAES.Config{
			Key: ctx.String("key"),
		},
	)
	return
}

func (a *aes) cipher(value string) (string, error) {
	log.Debug("Ciphering using AES")
	ciphered, err := a.client.Cipher(value)
	if err != nil {
		return "", err
	}
	return ciphered, nil
}

func (a *aes) decipher(value string) (string, error) {
	log.Debugf("Deciphering '%s' using AES", value)
	plain, err := a.client.Decipher(value)
	if err != nil {
		return "", err
	}
	return plain, nil
}
