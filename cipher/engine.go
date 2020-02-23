package cipher

import (
	"fmt"
	"regexp"

	"github.com/mvisonneau/s5/cipher/aes"
	"github.com/mvisonneau/s5/cipher/aws"
	"github.com/mvisonneau/s5/cipher/gcp"
	"github.com/mvisonneau/s5/cipher/pgp"
	"github.com/mvisonneau/s5/cipher/vault"
)

const (
	// InputRegexp is defining the syntax of an s5 input variable
	InputRegexp string = `{{\s?s5:([A-Za-z0-9+\\/=]*)\s?}}`
)

// Engine is an interface of supported/required commands for each cipher engine
type Engine interface {
	Cipher(string) (string, error)
	Decipher(string) (string, error)
}

// NewAESClient creates a AES client
func NewAESClient(key string) (*aes.Client, error) {
	return aes.NewClient(&aes.Config{
		Key: key,
	})
}

// NewAWSClient creates a AWS client
func NewAWSClient(kmsKeyArn string) (*aws.Client, error) {
	return aws.NewClient(&aws.Config{
		KmsKeyArn: kmsKeyArn,
	})
}

// NewGCPClient creates a GCP client
func NewGCPClient(kmsKeyName string) (*gcp.Client, error) {
	return gcp.NewClient(&gcp.Config{
		KmsKeyName: kmsKeyName,
	})
}

// NewPGPClient creates a PGP client
func NewPGPClient(publicKeyPath, privateKeyPath string) (*pgp.Client, error) {
	return pgp.NewClient(&pgp.Config{
		PublicKeyPath:  publicKeyPath,
		PrivateKeyPath: privateKeyPath,
	})
}

// NewVaultClient creates a Vault client
func NewVaultClient(key string) (*vault.Client, error) {
	return vault.NewClient(&vault.Config{
		Key: key,
	})
}

// GenerateOutput return a ciphered string in a s5 format
func GenerateOutput(value string) string {
	return fmt.Sprintf("{{s5:%s}}", value)
}

// ParseInput retrieves ciphered value from a string in the s5 format
func ParseInput(value string) (string, error) {
	re := regexp.MustCompile(InputRegexp)
	if !re.MatchString(value) {
		return "", fmt.Errorf("Invalid string format, should be '{{s5:*}}'")
	}
	return re.FindStringSubmatch(value)[1], nil
}
