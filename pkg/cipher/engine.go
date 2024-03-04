package cipher

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"

	"github.com/mvisonneau/s5/pkg/cipher/aes"
	"github.com/mvisonneau/s5/pkg/cipher/aws"
	"github.com/mvisonneau/s5/pkg/cipher/gcp"
	"github.com/mvisonneau/s5/pkg/cipher/pgp"
	"github.com/mvisonneau/s5/pkg/cipher/vault"
)

const (
	// InputRegexp is defining the syntax of an s5 input variable.
	InputRegexp string = `{{\s?s5:([A-Za-z0-9+\\/=]*)\s?}}`
)

// Engine is an interface of supported/required commands for each cipher engine.
type Engine interface {
	Cipher(value string) (string, error)
	Decipher(value string) (string, error)
}

// NewAESClient creates a AES client.
func NewAESClient(key string) (*aes.Client, error) {
	c, err := aes.NewClient(&aes.Config{
		Key: key,
	})
	if err != nil {
		return c, errors.Wrap(err, "creating new AES engine client")
	}

	return c, nil
}

// NewAWSClient creates a AWS client.
func NewAWSClient(kmsKeyArn string) (*aws.Client, error) {
	c, err := aws.NewClient(&aws.Config{
		KmsKeyArn: kmsKeyArn,
	})
	if err != nil {
		return c, errors.Wrap(err, "creating new AWS engine client")
	}

	return c, nil
}

// NewGCPClient creates a GCP client.
func NewGCPClient(kmsKeyName string) (*gcp.Client, error) {
	c, err := gcp.NewClient(&gcp.Config{
		KmsKeyName: kmsKeyName,
	})
	if err != nil {
		return c, errors.Wrap(err, "creating new GCP engine client")
	}

	return c, nil
}

// NewPGPClient creates a PGP client.
func NewPGPClient(publicKeyPath, privateKeyPath string) (*pgp.Client, error) {
	c, err := pgp.NewClient(&pgp.Config{
		PublicKeyPath:  publicKeyPath,
		PrivateKeyPath: privateKeyPath,
	})
	if err != nil {
		return c, errors.Wrap(err, "creating new PGP engine client")
	}

	return c, nil
}

// NewVaultClient creates a Vault client.
func NewVaultClient(key string) (*vault.Client, error) {
	c, err := vault.NewClient(&vault.Config{
		Key: key,
	})
	if err != nil {
		return c, errors.Wrap(err, "creating new Vault engine client")
	}

	return c, nil
}

// GenerateOutput return a ciphered string in a s5 format.
func GenerateOutput(value string) string {
	return fmt.Sprintf("{{s5:%s}}", value)
}

// ParseInput retrieves ciphered value from a string in the s5 format.
func ParseInput(value string) (string, error) {
	re := regexp.MustCompile(InputRegexp)
	if !re.MatchString(value) {
		return "", errors.New("invalid string format, should be '{{s5:*}}'")
	}

	return re.FindStringSubmatch(value)[1], nil
}
