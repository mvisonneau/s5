package pgp

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/pkg/errors"

	"github.com/mvisonneau/s5/internal/logs"
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

	publicKey, err = os.ReadFile(config.PublicKeyPath)
	if err != nil {
		return nil, errors.New("reading the public-key file")
	}

	var entity *openpgp.Entity
	if len(config.PrivateKeyPath) > 0 {
		privateKey, err = os.ReadFile(config.PrivateKeyPath)
		if err != nil {
			return nil, errors.New("reading the private-key file")
		}

		entityList, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(privateKey))
		if err != nil || len(entityList) == 0 {
			return nil, errors.New("parsing the private key")
		}
		entity = entityList[0]
	} else {
		entityList, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(publicKey))
		if err != nil || len(entityList) == 0 {
			return nil, errors.New("parsing the public key")
		}
		entity = entityList[0]
	}

	return &Client{Entity: entity}, nil
}

// Cipher a value using a Public pgp key.
func (c *Client) Cipher(ctx context.Context, value string) (string, error) {
	logs.LoggerFromContext(ctx).Debug("ciphering using a public pgp key")
	var buf bytes.Buffer
	w, err := armor.Encode(&buf, "PGP MESSAGE", nil)
	if err != nil {
		return "", errors.Wrap(err, "armor encode")
	}

	plaintextWriter, err := openpgp.Encrypt(w, []*openpgp.Entity{c.Entity}, nil, nil, nil)
	if err != nil {
		return "", errors.Wrap(err, "encrypt")
	}

	_, err = plaintextWriter.Write([]byte(value))
	if err != nil {
		return "", errors.Wrap(err, "write to plaintext writer")
	}

	_ = plaintextWriter.Close()
	_ = w.Close()

	return buf.String(), nil
}

// Decipher a value using a Public and Private pgp keypair.
func (c *Client) Decipher(ctx context.Context, value string) (string, error) {
	logs.LoggerFromContext(ctx).Debug("deciphering using a public pgp key")
	block, err := armor.Decode(strings.NewReader(value))
	if err != nil {
		return "", errors.Wrap(err, "armor decode")
	}

	md, err := openpgp.ReadMessage(block.Body, openpgp.EntityList{c.Entity}, nil, nil)
	if err != nil {
		return "", errors.Wrap(err, "read message")
	}

	decryptedBytes, err := io.ReadAll(md.UnverifiedBody)
	if err != nil {
		return "", errors.Wrap(err, "read decrypted body")
	}

	return string(decryptedBytes), nil
}
