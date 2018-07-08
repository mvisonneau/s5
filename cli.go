package main

import (
	"github.com/urfave/cli"
)

var version = "<devel>"

// runCli : Generates cli configuration for the application
func runCli() (c *cli.App) {
	c = cli.NewApp()
	c.Name = "s5"
	c.Version = version
	c.Usage = "cipher/decipher text within a file from a (Hashicorp) Vault transit key"
	c.EnableBashCompletion = true

	c.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "vault-addr",
			EnvVar:      "VAULT_ADDR",
			Usage:       "vault `address`",
			Destination: &cfg.Vault.Address,
		},
		cli.StringFlag{
			Name:        "vault-token",
			EnvVar:      "VAULT_TOKEN",
			Usage:       "vault `token`",
			Destination: &cfg.Vault.Token,
		},
		cli.StringFlag{
			Name:        "transit-key,k",
			EnvVar:      "S5_TRANSIT_KEY",
			Usage:       "`name` of the transit key used by s5 to cipher/decipher data",
			Value:       "default",
			Destination: &cfg.TransitKey,
		},
		cli.StringFlag{
			Name:        "log-level",
			EnvVar:      "S5_LOG_LEVEL",
			Usage:       "log `level` (debug,info,warn,fatal,panic)",
			Value:       "info",
			Destination: &cfg.Log.Level,
		},
		cli.StringFlag{
			Name:        "log-format",
			EnvVar:      "S5_LOG_FORMAT",
			Usage:       "log `format` (json,text)",
			Value:       "text",
			Destination: &cfg.Log.Format,
		},
	}

	c.Commands = []cli.Command{
		{
			Name:      "cipher",
			Usage:     "return an encrypted s5 pattern that can be included in any file",
			ArgsUsage: "<value>",
			Action:    run,
		},
		{
			Name:      "decipher",
			Usage:     "return an unencrypted s5 value from a given pattern",
			ArgsUsage: "<pattern>",
			Action:    run,
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
			Action:    run,
		},
	}

	return
}
