package command

import (
	"fmt"

	cipherPGP "github.com/mvisonneau/s5/cipher/pgp"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type pgp struct {
	client *cipherPGP.Client
}

func (g *pgp) configure(ctx *cli.Context) (err error) {
	if len(ctx.String("public-key")) == 0 {
		return fmt.Errorf("You need to specify the public-key path (--public-key - $S5_pgp_PUBLIC_KEY_PATH)")
	}

	g.client, err = cipherPGP.Init(
		&cipherPGP.Config{
			PublicKeyPath:  ctx.String("public-key"),
			PrivateKeyPath: ctx.String("private-key"),
		},
	)
	return
}

func (g *pgp) cipher(value string) (string, error) {
	log.Debug("Ciphering using pgp public key")
	ciphered, err := g.client.Cipher(value)
	if err != nil {
		return "", err
	}
	return ciphered, nil
}

func (g *pgp) decipher(value string) (string, error) {
	log.Debugf("Deciphering '%s' using pgp public/private keypair", value)
	plain, err := g.client.Decipher(value)
	if err != nil {
		return "", err
	}
	return plain, nil
}
