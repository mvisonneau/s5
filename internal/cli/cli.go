package cli

import (
	"github.com/urfave/cli/v3"

	"github.com/mvisonneau/s5/internal/app"
	"github.com/mvisonneau/s5/internal/cmd"
)

// NewApp configures the CLI application.
func NewApp() *cli.Command {
	return &cli.Command{
		Name:                  app.Name,
		Version:               app.Version,
		Usage:                 "cipher/decipher text within a file",
		HideHelpCommand:       true,
		EnableShellCompletion: true,

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-level",
				Sources: cli.EnvVars("S5_LOG_LEVEL"),
				Usage:   "log `level` (debug,info,warn,fatal,panic)",
				Value:   "info",
			},
			&cli.StringFlag{
				Name:    "log-format",
				Sources: cli.EnvVars("S5_LOG_FORMAT"),
				Usage:   "log `format` (json,text)",
				Value:   "text",
			},
		},

		Commands: []*cli.Command{
			{
				Name:  "cipher",
				Usage: "return an encrypted s5 pattern that can be included in any file",
				Commands: []*cli.Command{
					{
						Name:      "aes",
						Usage:     "cipher using an AES key",
						ArgsUsage: "<value>",
						Action:    cmd.Execute(cmd.Cipher),
						Flags:     append(aesFlags, noTrimFlag),
					},
					{
						Name:      "aws",
						Usage:     "cipher using an AWS KMS key",
						ArgsUsage: "<value>",
						Action:    cmd.Execute(cmd.Cipher),
						Flags:     append(awsFlags, noTrimFlag),
					},
					{
						Name:      "gcp",
						Usage:     "cipher using a GCP KMS key",
						ArgsUsage: "<value>",
						Action:    cmd.Execute(cmd.Cipher),
						Flags:     append(gcpFlags, noTrimFlag),
					},
					{
						Name:      "pgp",
						Usage:     "cipher using a public pgp key",
						ArgsUsage: "<value>",
						Action:    cmd.Execute(cmd.Cipher),
						Flags: cli.FlagsByName{
							noTrimFlag,
							pgpPublicKeyPathFlag,
						},
					},
					{
						Name:      "vault",
						Usage:     "cipher using a vault transit key",
						ArgsUsage: "<value>",
						Action:    cmd.Execute(cmd.Cipher),
						Flags:     append(vaultFlags, noTrimFlag),
					},
				},
			},
			{
				Name:  "decipher",
				Usage: "return an unencrypted s5 value from a given pattern",
				Commands: []*cli.Command{
					{
						Name:      "aes",
						Usage:     "decipher using an AES key",
						ArgsUsage: "<value>",
						Action:    cmd.Execute(cmd.Decipher),
						Flags:     aesFlags,
					},
					{
						Name:      "aws",
						Usage:     "decipher using an AWS KMS key",
						ArgsUsage: "<value>",
						Action:    cmd.Execute(cmd.Decipher),
					},
					{
						Name:      "gcp",
						Usage:     "decipher using a GCP key",
						ArgsUsage: "<value>",
						Action:    cmd.Execute(cmd.Decipher),
						Flags:     gcpFlags,
					},
					{
						Name:      "pgp",
						Usage:     "decipher using a public/private pgp keypair",
						ArgsUsage: "<value>",
						Action:    cmd.Execute(cmd.Decipher),
						Flags: cli.FlagsByName{
							pgpPublicKeyPathFlag,
							pgpPrivateKeyPathFlag,
						},
					},
					{
						Name:      "vault",
						Usage:     "decipher using a vault transit key",
						ArgsUsage: "<value>",
						Action:    cmd.Execute(cmd.Decipher),
						Flags:     vaultFlags,
					},
				},
			},
			{
				Name:  "render",
				Usage: "render a file that (may) contain s5 encrypted patterns",
				Commands: []*cli.Command{
					{
						Name:      "aes",
						Usage:     "render using an AES key",
						ArgsUsage: "<value>",
						Action:    cmd.Execute(cmd.Render),
						Flags:     append(renderFlags, aesFlags...),
					},
					{
						Name:      "aws",
						Usage:     "render using an AWS KMS key",
						ArgsUsage: "<value>",
						Action:    cmd.Execute(cmd.Render),
						Flags:     renderFlags,
					},
					{
						Name:      "gcp",
						Usage:     "render using a GCP key",
						ArgsUsage: "<value>",
						Action:    cmd.Execute(cmd.Render),
						Flags:     append(renderFlags, gcpFlags...),
					},
					{
						Name:      "pgp",
						Usage:     "render using a public/private pgp keypair",
						ArgsUsage: "<value>",
						Action:    cmd.Execute(cmd.Render),
						Flags:     append(renderFlags, pgpPublicKeyPathFlag, pgpPrivateKeyPathFlag),
					},
					{
						Name:      "vault",
						Usage:     "render using a vault transit key",
						ArgsUsage: "<filename> [--in-place] [--output <output_file>]",
						Action:    cmd.Execute(cmd.Render),
						Flags:     append(renderFlags, vaultFlags...),
					},
				},
			},
		},
	}
}
