package cmd

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/mvisonneau/s5/pkg/cipher"
)

// Decipher is used for the decipher commands.
func Decipher(_ context.Context, cmd *cli.Command) error {
	cipherEngine, err := getCipherEngine(cmd)
	if err != nil {
		return err
	}

	if err := configure(cmd); err != nil {
		return err
	}

	input, err := readInput(cmd)
	if err != nil {
		if err := cli.ShowSubcommandHelp(cmd); err != nil {
			return errors.Wrap(err, "rendering subcommand help")
		}

		return err
	}

	log.Debug("Validating input string")

	parsedInput, err := cipher.ParseInput(input)
	if err != nil {
		return errors.Wrap(err, "parsing input")
	}

	deciphered, err := cipherEngine.Decipher(parsedInput)
	if err != nil {
		return errors.Wrap(err, "deciphering input")
	}

	fmt.Print(deciphered)

	return nil
}
