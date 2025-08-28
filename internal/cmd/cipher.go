package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/mvisonneau/s5/pkg/cipher"
)

// Cipher is used for the cipher commands.
func Cipher(ctx context.Context, cmd *cli.Command) error {
	cipherEngine, err := getCipherEngine(ctx, cmd)
	if err != nil {
		return err
	}

	input, err := readInput(cmd)
	if err != nil {
		if err = cli.ShowSubcommandHelp(cmd); err != nil {
			return errors.Wrap(err, "rendering subcommand help")
		}

		return err
	}

	if !cmd.Bool("no-trim") {
		input = strings.Trim(input, " \n")
	}

	ciphered, err := cipherEngine.Cipher(ctx, input)
	if err != nil {
		return errors.Wrap(err, "ciphering input")
	}

	fmt.Print(cipher.GenerateOutput(ciphered))

	return nil
}
