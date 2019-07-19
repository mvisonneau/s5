package cli

import (
	"time"

	"github.com/mvisonneau/s5/command"
	"github.com/urfave/cli"
)

// Init : Generates CLI configuration for the application
func Init(version *string, start time.Time) (app *cli.App) {
	app = cli.NewApp()
	app.Name = "s5"
	app.Version = *version
	app.Usage = "cipher/decipher text within a file"
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "log-level",
			EnvVar: "S5_LOG_LEVEL",
			Usage:  "log `level` (debug,info,warn,fatal,panic)",
			Value:  "info",
		},
		cli.StringFlag{
			Name:   "log-format",
			EnvVar: "S5_LOG_FORMAT",
			Usage:  "log `format` (json,text)",
			Value:  "text",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "cipher",
			Usage: "return an encrypted s5 pattern that can be included in any file",
			Subcommands: []cli.Command{
				{
					Name:      "aes",
					Usage:     "cipher using AES",
					ArgsUsage: "<value>",
					Action:    command.Cipher,
					Flags:     aesFlags,
				},
				{
					Name:      "pgp",
					Usage:     "cipher using a public pgp key",
					ArgsUsage: "<value>",
					Action:    command.Cipher,
					Flags: []cli.Flag{
						pgpPublicKeyPathFlag,
					},
				},
				{
					Name:      "vault",
					Usage:     "cipher using a vault transit key",
					ArgsUsage: "<value>",
					Action:    command.Cipher,
					Flags:     vaultFlags,
				},
			},
		},
		{
			Name:  "decipher",
			Usage: "return an unencrypted s5 value from a given pattern",
			Subcommands: []cli.Command{
				{
					Name:      "aes",
					Usage:     "decipher using AES",
					ArgsUsage: "<value>",
					Action:    command.Decipher,
					Flags:     aesFlags,
				},
				{
					Name:      "pgp",
					Usage:     "decipher using a public/private pgp keypair",
					ArgsUsage: "<value>",
					Action:    command.Decipher,
					Flags: []cli.Flag{
						pgpPublicKeyPathFlag,
						pgpPrivateKeyPathFlag,
					},
				},
				{
					Name:      "vault",
					Usage:     "decipher using a vault transit key",
					ArgsUsage: "<value>",
					Action:    command.Decipher,
					Flags:     vaultFlags,
				},
			},
		},
		{
			Name:  "render",
			Usage: "render a file that (may) contain s5 encrypted patterns",
			Subcommands: []cli.Command{
				{
					Name:      "aes",
					Usage:     "render using AES",
					ArgsUsage: "<value>",
					Action:    command.Render,
					Flags:     append(renderFlags, aesFlags...),
				},
				{
					Name:      "pgp",
					Usage:     "render using a public/private pgp keypair",
					ArgsUsage: "<value>",
					Action:    command.Render,
					Flags:     append(renderFlags, pgpPublicKeyPathFlag, pgpPrivateKeyPathFlag),
				},
				{
					Name:      "vault",
					Usage:     "render using a vault transit key",
					ArgsUsage: "<filename> [--in-place] [--output <output_file>]",
					Action:    command.Render,
					Flags:     append(renderFlags, vaultFlags...),
				},
			},
		},
	}

	app.Metadata = map[string]interface{}{
		"startTime": start,
	}

	return
}
