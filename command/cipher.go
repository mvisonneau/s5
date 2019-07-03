package command

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// Cipher function
func Cipher(ctx *cli.Context) error {
	if err := configure(ctx); err != nil {
		return cli.NewExitError(err, 1)
	}

	input, err := readInput(ctx)
	if err != nil {
		cli.ShowSubcommandHelp(ctx)
		return exit(err, 1)
	}

	log.Debug("Ciphering using Vault transit key")
	ciphered, err := vt.Cipher(input)
	if err != nil {
		return exit(err, 1)
	}
	fmt.Printf("{{ %s }}", ciphered)

	return exit(nil, 0)
}
