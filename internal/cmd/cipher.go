package cmd

import (
	"fmt"
	"strings"

	"github.com/mvisonneau/s5/pkg/cipher"
	"github.com/urfave/cli/v2"
)

// Cipher is used for the cipher commands
func Cipher(ctx *cli.Context) (int, error) {
	cipherEngine, err := getCipherEngine(ctx)
	if err != nil {
		return 1, err
	}

	if err := configure(ctx); err != nil {
		return 1, err
	}

	input, err := readInput(ctx)
	if err != nil {
		if err = cli.ShowSubcommandHelp(ctx); err != nil {
			return 1, err
		}
		return 1, err
	}

	if !ctx.Bool("no-trim") {
		input = strings.Trim(input, " \n")
	}

	ciphered, err := cipherEngine.Cipher(input)
	if err != nil {
		return 1, err
	}

	fmt.Print(cipher.GenerateOutput(ciphered))

	return 0, nil
}
