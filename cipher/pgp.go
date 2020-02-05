package cipher

import (
	"fmt"

	cipherPGP "github.com/mvisonneau/s5/cipher/pgp"

	log "github.com/sirupsen/logrus"
)

type PGP struct {
	client *cipherPGP.Client
}

func NewPGP(publicKeyPath, privateKeyPath string) (*PGP, error) {
	if len(publicKeyPath) == 0 {
		return nil, fmt.Errorf("You need to specify the public-key path for GPG")
	}

	c, err := cipherPGP.Init(
		&cipherPGP.Config{
			PublicKeyPath:  publicKeyPath,
			PrivateKeyPath: privateKeyPath,
		},
	)

	if err != nil {
		return nil, err
	}

	return &PGP{c}, nil
}

func (g *PGP) Cipher(value string) (string, error) {
	log.Debug("Ciphering using pgp public key")
	ciphered, err := g.client.Cipher(value)
	if err != nil {
		return "", err
	}
	return ciphered, nil
}

func (g *PGP) Decipher(value string) (string, error) {
	log.Debugf("Deciphering '%s' using pgp public/private keypair", value)
	plain, err := g.client.Decipher(value)
	if err != nil {
		return "", err
	}
	return plain, nil
}
