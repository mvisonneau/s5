package command

import (
	cipherVault "github.com/mvisonneau/s5/cipher/vault"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type vault struct {
	client *cipherVault.Client
}

func (v *vault) configure(ctx *cli.Context) (err error) {
	v.client, err = cipherVault.Init(
		&cipherVault.Config{
			Key: ctx.String("transit-key"),
		},
	)
	return
}

func (v *vault) cipher(value string) (string, error) {
	log.Debug("Ciphering using Vault transit key")
	ciphered, err := v.client.Cipher(value)
	if err != nil {
		return "", err
	}
	return ciphered, nil
}

func (v *vault) decipher(value string) (string, error) {
	log.Debugf("Deciphering '%s' using Vault transit key", value)
	plain, err := v.client.Decipher(value)
	if err != nil {
		return "", err
	}
	return plain, nil
}
