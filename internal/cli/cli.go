package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/mvisonneau/s5/internal/cmd"
	"github.com/urfave/cli/v2"
)

// Run handles the instanciation of the CLI application
func Run(version string, args []string) {
	err := NewApp(version, time.Now()).Run(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// NewApp configures the CLI application
func NewApp(version string, start time.Time) (app *cli.App) {
	app = cli.NewApp()
	app.Name = "s5"
	app.Version = version
	app.Usage = "cipher/decipher text within a file"
	app.EnableBashCompletion = true

	app.Flags = cli.FlagsByName{
		&cli.StringFlag{
			Name:    "log-level",
			EnvVars: []string{"S5_LOG_LEVEL"},
			Usage:   "log `level` (debug,info,warn,fatal,panic)",
			Value:   "info",
		},
		&cli.StringFlag{
			Name:    "log-format",
			EnvVars: []string{"S5_LOG_FORMAT"},
			Usage:   "log `format` (json,text)",
			Value:   "text",
		},
	}

	app.Commands = cli.CommandsByName{
		{
			Name:  "cipher",
			Usage: "return an encrypted s5 pattern that can be included in any file",
			Subcommands: cli.CommandsByName{
				{
					Name:      "aes",
					Usage:     "cipher using an AES key",
					ArgsUsage: "<value>",
					Action:    cmd.ExecWrapper(cmd.Cipher),
					Flags:     append(aesFlags, noTrimFlag),
				},
				{
					Name:      "aws",
					Usage:     "cipher using an AWS KMS key",
					ArgsUsage: "<value>",
					Action:    cmd.ExecWrapper(cmd.Cipher),
					Flags:     append(awsFlags, noTrimFlag),
				},
				{
					Name:      "gcp",
					Usage:     "cipher using a GCP KMS key",
					ArgsUsage: "<value>",
					Action:    cmd.ExecWrapper(cmd.Cipher),
					Flags:     append(gcpFlags, noTrimFlag),
				},
				{
					Name:      "pgp",
					Usage:     "cipher using a public pgp key",
					ArgsUsage: "<value>",
					Action:    cmd.ExecWrapper(cmd.Cipher),
					Flags: cli.FlagsByName{
						noTrimFlag,
						pgpPublicKeyPathFlag,
					},
				},
				{
					Name:      "vault",
					Usage:     "cipher using a vault transit key",
					ArgsUsage: "<value>",
					Action:    cmd.ExecWrapper(cmd.Cipher),
					Flags:     append(vaultFlags, noTrimFlag),
				},
			},
		},
		{
			Name:  "decipher",
			Usage: "return an unencrypted s5 value from a given pattern",
			Subcommands: cli.CommandsByName{
				{
					Name:      "aes",
					Usage:     "decipher using an AES key",
					ArgsUsage: "<value>",
					Action:    cmd.ExecWrapper(cmd.Decipher),
					Flags:     aesFlags,
				},
				{
					Name:      "aws",
					Usage:     "decipher using an AWS KMS key",
					ArgsUsage: "<value>",
					Action:    cmd.ExecWrapper(cmd.Decipher),
				},
				{
					Name:      "gcp",
					Usage:     "decipher using a GCP key",
					ArgsUsage: "<value>",
					Action:    cmd.ExecWrapper(cmd.Decipher),
					Flags:     gcpFlags,
				},
				{
					Name:      "pgp",
					Usage:     "decipher using a public/private pgp keypair",
					ArgsUsage: "<value>",
					Action:    cmd.ExecWrapper(cmd.Decipher),
					Flags: cli.FlagsByName{
						pgpPublicKeyPathFlag,
						pgpPrivateKeyPathFlag,
					},
				},
				{
					Name:      "vault",
					Usage:     "decipher using a vault transit key",
					ArgsUsage: "<value>",
					Action:    cmd.ExecWrapper(cmd.Decipher),
					Flags:     vaultFlags,
				},
			},
		},
		{
			Name:  "render",
			Usage: "render a file that (may) contain s5 encrypted patterns",
			Subcommands: cli.CommandsByName{
				{
					Name:      "aes",
					Usage:     "render using an AES key",
					ArgsUsage: "<value>",
					Action:    cmd.ExecWrapper(cmd.Render),
					Flags:     append(renderFlags, aesFlags...),
				},
				{
					Name:      "aws",
					Usage:     "render using an AWS KMS key",
					ArgsUsage: "<value>",
					Action:    cmd.ExecWrapper(cmd.Render),
					Flags:     renderFlags,
				},
				{
					Name:      "gcp",
					Usage:     "render using a GCP key",
					ArgsUsage: "<value>",
					Action:    cmd.ExecWrapper(cmd.Render),
					Flags:     append(renderFlags, gcpFlags...),
				},
				{
					Name:      "pgp",
					Usage:     "render using a public/private pgp keypair",
					ArgsUsage: "<value>",
					Action:    cmd.ExecWrapper(cmd.Render),
					Flags:     append(renderFlags, pgpPublicKeyPathFlag, pgpPrivateKeyPathFlag),
				},
				{
					Name:      "vault",
					Usage:     "render using a vault transit key",
					ArgsUsage: "<filename> [--in-place] [--output <output_file>]",
					Action:    cmd.ExecWrapper(cmd.Render),
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
