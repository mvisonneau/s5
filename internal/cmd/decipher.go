package cmd

import (
	"fmt"

	"github.com/mvisonneau/s5/pkg/cipher"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Decipher is used for the decipher commands.
func Decipher(ctx *cli.Context) (int, error) {
	cipherEngine, err := getCipherEngine(ctx)
	if err != nil {
		return 1, err
	}

	if err := configure(ctx); err != nil {
		return 1, err
	}

	input, err := readInput(ctx)
	if err != nil {
		if err := cli.ShowSubcommandHelp(ctx); err != nil {
			return 1, errors.Wrap(err, "rendering subcommand help")
		}

		return 1, err
	}

	log.Debug("Validating input string")

	parsedInput, err := cipher.ParseInput(input)
	if err != nil {
		return 1, errors.Wrap(err, "parsing input")
	}

	deciphered, err := cipherEngine.Decipher(parsedInput)
	if err != nil {
		return 1, errors.Wrap(err, "deciphering input")
	}

	fmt.Print(deciphered) //nolint

	return 0, nil
}
