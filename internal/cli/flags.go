package cli

import (
	"github.com/urfave/cli/v3"
)

var aesFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "key",
		Sources: cli.EnvVars("S5_AES_KEY"),
		Usage:   "`path` to a readable key for AES encryption/decryption",
	},
}

var awsFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "kms-key-arn",
		Sources: cli.EnvVars("S5_AWS_KMS_KEY_ARN"),
		Usage:   "`arn` of a usable AWS KMS key for (de)ciphering text",
	},
}

var gcpFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "kms-key-name",
		Sources: cli.EnvVars("S5_GCP_KMS_KEY_NAME"),
		Usage:   "`name` of a usable GCP KMS key for (de)ciphering text",
	},
}

var noTrimFlag = &cli.BoolFlag{
	Name:    "no-trim",
	Sources: cli.EnvVars("S5_NO_TRIM"),
	Usage:   "if set, will not trim the input string from whitespaces and return carriages",
}

var pgpPublicKeyPathFlag = &cli.StringFlag{
	Name:    "public-key-path",
	Sources: cli.EnvVars("S5_PGP_PUBLIC_KEY_PATH"),
	Usage:   "`path` to a readable public pgp key (armored)",
}

var pgpPrivateKeyPathFlag = &cli.StringFlag{
	Name:    "private-key-path",
	Sources: cli.EnvVars("S5_PGP_PRIVATE_KEY_PATH"),
	Usage:   "`path` to a readable private pgp key (armored)",
}

var renderFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Usage:   "output `filename`",
	},
	&cli.BoolFlag{
		Name:    "in-place",
		Aliases: []string{"i"},
		Usage:   " ",
	},
}

var vaultFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "transit-key",
		Sources: cli.EnvVars("S5_VAULT_TRANSIT_KEY"),
		Usage:   "`name` of the transit key used by s5 to cipher/decipher data",
		Value:   "default",
	},
}
