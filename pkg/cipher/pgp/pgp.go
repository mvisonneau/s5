package pgp

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/jchavannes/go-pgp/pgp"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/openpgp"
)

// Config handles necessary information for PGP.
type Config struct {
	PublicKeyPath  string
	PrivateKeyPath string
}

// Client can be used to execute ciphering actions.
type Client struct {
	Entity *openpgp.Entity
}

// NewClient reads PGP key values and return a entity through a client object.
func NewClient(config *Config) (*Client, error) {
	var (
		err                   error
		publicKey, privateKey []byte
	)

	publicKey, err = ioutil.ReadFile(config.PublicKeyPath)
	if err != nil {
		return nil, errors.New("reading the public-key file")
	}

	if len(config.PrivateKeyPath) > 0 {
		privateKey, err = ioutil.ReadFile(config.PrivateKeyPath)
		if err != nil {
			return nil, errors.New("reading the private-key file")
		}
	}

	var pgpEntity *openpgp.Entity

	pgpEntity, err = pgp.GetEntity(publicKey, privateKey)
	if err != nil {
		log.Debugf("PUBLIC KEY:\n%s", string(publicKey))
		log.Debugf("PRIVATE KEY:\n%s", string(privateKey))

		return nil, errors.New("creating the PGP entity from the key(s)")
	}

	return &Client{Entity: pgpEntity}, nil
}

// Cipher a value using a Public pgp key.
func (c *Client) Cipher(value string) (string, error) {
	log.Debug("Ciphering using pgp public key")

	d, err := pgp.Encrypt(c.Entity, []byte(value))
	if err != nil {
		return "", errors.Wrap(err, "ciphering using PGP")
	}

	return base64.StdEncoding.EncodeToString(d), nil
}

// Decipher a value using a Public and Private pgp keypair.
func (c *Client) Decipher(value string) (string, error) {
	log.Debugf("Deciphering '%s' using pgp public/private keypair", value)

	str, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("base64decode - value : '%s'", value))
	}

	d, err := pgp.Decrypt(c.Entity, str)
	if err != nil {
		return "", errors.Wrap(err, "deciphering using PGP")
	}

	return string(d), nil
}
