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
	app.Usage = "cipher/decipher text within a file from a (Hashicorp) Vault transit key"
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "vault-transit-key",
			EnvVar: "S5_VAULT_TRANSIT_KEY",
			Usage:  "`name` of the transit key used by s5 to cipher/decipher data",
			Value:  "default",
		},
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
			Name:      "cipher",
			Usage:     "return an encrypted s5 pattern that can be included in any file",
			ArgsUsage: "<value>",
			Action:    command.Cipher,
		},
		{
			Name:      "decipher",
			Usage:     "return an unencrypted s5 value from a given pattern",
			ArgsUsage: "<pattern>",
			Action:    command.Decipher,
		},
		{
			Name:      "render",
			Usage:     "render a file that (may) contain s5 encrypted patterns",
			ArgsUsage: "<filename> [--in-place] [--output <output_file>]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output,o",
					Usage: "output `filename`",
				},
				cli.BoolFlag{
					Name:  "in-place,i",
					Usage: " ",
				},
			},
			Action: command.Render,
		},
	}

	app.Metadata = map[string]interface{}{
		"startTime": start,
	}

	return
}
