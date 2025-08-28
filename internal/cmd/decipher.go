package cmd

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/mvisonneau/s5/pkg/cipher"
)

// Decipher is used for the decipher commands.
func Decipher(ctx context.Context, cmd *cli.Command) error {
	cipherEngine, err := getCipherEngine(ctx, cmd)
	if err != nil {
		return err
	}

	input, err := readInput(cmd)
	if err != nil {
		if err := cli.ShowSubcommandHelp(cmd); err != nil {
			return errors.Wrap(err, "rendering subcommand help")
		}

		return err
	}

	parsedInput, err := cipher.ParseInput(input)
	if err != nil {
		return errors.Wrap(err, "parsing input")
	}

	deciphered, err := cipherEngine.Decipher(ctx, parsedInput)
	if err != nil {
		return errors.Wrap(err, "deciphering input")
	}

	fmt.Print(deciphered)

	return nil
}
