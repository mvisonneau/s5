package pgp

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/jchavannes/go-pgp/pgp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/openpgp"
)

// Config handles necessary information for PGP
type Config struct {
	PublicKeyPath  string
	PrivateKeyPath string
}

// Client can be used to execute ciphering actions
type Client struct {
	Entity *openpgp.Entity
}

// Init : Reads PGP key values and return a entity
func Init(config *Config) (*Client, error) {
	var err error
	publicKey, privateKey := []byte{}, []byte{}

	publicKey, err = ioutil.ReadFile(config.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("Error while reading the public-key file")
	}

	if len(config.PrivateKeyPath) > 0 {
		privateKey, err = ioutil.ReadFile(config.PrivateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("Error while reading the private-key file")
		}
	}

	var pgpEntity *openpgp.Entity
	pgpEntity, err = pgp.GetEntity(publicKey, privateKey)
	if err != nil {
		log.Debugf("PUBLIC KEY:\n%s", string(publicKey))
		log.Debugf("PRIVATE KEY:\n%s", string(privateKey))
		return nil, fmt.Errorf("Error while creating the PGP entity from the key(s)")
	}

	return &Client{Entity: pgpEntity}, nil
}

// Cipher a value using a Public pgp key
func (c *Client) Cipher(value string) (string, error) {
	d, err := pgp.Encrypt(c.Entity, []byte(value))
	if err != nil {
		return "", fmt.Errorf("PGP error : %s", err)
	}

	return base64.StdEncoding.EncodeToString(d), nil
}

// Decipher a value using a Public and Private pgp keypair
func (c *Client) Decipher(value string) (string, error) {
	str, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", fmt.Errorf("base64decode error : %s - value : %s", err, value)
	}

	d, err := pgp.Decrypt(c.Entity, str)
	if err != nil {
		return "", fmt.Errorf("PGP error : %s", err)
	}

	return string(d), nil
}
