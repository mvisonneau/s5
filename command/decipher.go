package command

import (
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// Decipher function
func Decipher(ctx *cli.Context) error {
	if err := configure(ctx); err != nil {
		return cli.NewExitError(err, 1)
	}

	input, err := readInput(ctx)
	if err != nil {
		cli.ShowSubcommandHelp(ctx)
		return exit(err, 1)
	}

	log.Debug("Validating input string")
	re := regexp.MustCompile("{{ (s5:.*) }}")
	if !re.MatchString(input) {
		return exit(fmt.Errorf("Invalid string format, should be '{{ s5:* }}'"), 1)
	}

	log.Debugf("Deciphering '%s' using Vault transit key", re.FindStringSubmatch(input)[1])
	plain, err := vt.Decipher(re.FindStringSubmatch(input)[1])
	if err != nil {
		return exit(err, 1)
	}
	fmt.Print(plain)

	return exit(nil, 0)
}
