package command

import (
	cipherGCP "github.com/mvisonneau/s5/cipher/gcp"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type gcp struct {
	client *cipherGCP.Client
}

func (g *gcp) configure(ctx *cli.Context) (err error) {
	g.client, err = cipherGCP.Init(
		&cipherGCP.Config{
			KmsKeyName: ctx.String("kms-key-name"),
		},
	)
	return
}

func (g *gcp) cipher(value string) (string, error) {
	log.Debug("Ciphering using a GCP KMS key")
	ciphered, err := g.client.Cipher(value)
	if err != nil {
		return "", err
	}
	return ciphered, nil
}

func (g *gcp) decipher(value string) (string, error) {
	log.Debugf("Deciphering '%s' using a GCP KMS key", value)
	plain, err := g.client.Decipher(value)
	if err != nil {
		return "", err
	}
	return plain, nil
}
