package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/hashicorp/vault/sdk/helper/mlock"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/mvisonneau/go-helpers/logger"
	"github.com/mvisonneau/s5/pkg/cipher"
)

var start time.Time

func configure(ctx *cli.Context) error {
	start = ctx.App.Metadata["startTime"].(time.Time)

	if err := logger.Configure(logger.Config{
		Level:  ctx.String("log-level"),
		Format: ctx.String("log-format"),
	}); err != nil {
		return errors.Wrap(err, "configuring logger")
	}

	return nil
}

func getCipherEngine(ctx *cli.Context) (engine cipher.Engine, err error) {
	cmds := ctx.Command.Names()
	switch cmds[len(cmds)-1] {
	case "aes":
		engine, err = cipher.NewAESClient(ctx.String("key"))
	case "aws":
		engine, err = cipher.NewAWSClient(ctx.String("kms-key-arn"))
	case "gcp":
		engine, err = cipher.NewGCPClient(ctx.String("kms-key-name"))
	case "pgp":
		engine, err = cipher.NewPGPClient(ctx.String("public-key-path"), ctx.String("private-key-path"))
	case "vault":
		engine, err = cipher.NewVaultClient(ctx.String("transit-key"))
	default:
		err = fmt.Errorf("engine %v is not implemented yet", ctx.Command.FullName())
	}

	if err != nil {
		return nil, errors.Wrap(err, "getting cipher engine")
	}

	return
}

func readInput(ctx *cli.Context) (input string, err error) {
	switch ctx.NArg() {
	case 0:
		read, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return "", errors.Wrap(err, "reading from stdin")
		}

		input = string(read)
	case 1:
		input = ctx.Args().First()
	default:
		return "", errors.New("too many arguments provided")
	}

	return
}

func exit(exitCode int, err error) cli.ExitCoder {
	defer log.WithFields(
		log.Fields{
			"execution-time": time.Since(start),
		},
	).Debug("exited..")

	if err != nil {
		log.Error(err.Error())
	}

	return cli.Exit("", exitCode)
}

// ExecWrapper gracefully logs and exits our `run` functions.
func ExecWrapper(f func(ctx *cli.Context) (int, error)) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		if err := mlock.LockMemory(); err != nil {
			log.WithError(err).Warn("s5 requires the IPC_LOCK capability in order to secure its memory")
		}

		return exit(f(ctx))
	}
}
