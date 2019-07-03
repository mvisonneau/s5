package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	vaultTransit "github.com/mvisonneau/s5/cipher/vault/transit"
	"github.com/mvisonneau/s5/logger"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var start time.Time
var vt *vaultTransit.Client

func configure(ctx *cli.Context) error {
	var err error
	start = ctx.App.Metadata["startTime"].(time.Time)

	lc := &logger.Config{
		Level:  ctx.GlobalString("log-level"),
		Format: ctx.GlobalString("log-format"),
	}

	if err = lc.Configure(); err != nil {
		return err
	}

	vt, err = vaultTransit.Init(
		&vaultTransit.Config{
			Key: ctx.String("transit-key"),
		},
	)

	return err
}

func readInput(c *cli.Context) (input string, err error) {
	switch c.NArg() {
	case 0:
		read, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		input = string(read)
	case 1:
		input = c.Args().First()
	default:
		return "", fmt.Errorf("Too many arguments provided")
	}
	return
}

func exit(err error, exitCode int) *cli.ExitError {
	defer log.Debugf("Executed in %s, exiting..", time.Since(start))
	if err != nil {
		log.Error(err.Error())
		return cli.NewExitError("", exitCode)
	}

	return cli.NewExitError("", 0)
}
