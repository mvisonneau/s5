package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/hashicorp/vault/sdk/helper/mlock"
	"github.com/mvisonneau/go-helpers/logger"
	"github.com/mvisonneau/s5/pkg/cipher"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var start time.Time

func configure(ctx *cli.Context) error {
	start = ctx.App.Metadata["startTime"].(time.Time)

	return logger.Configure(logger.Config{
		Level:  ctx.String("log-level"),
		Format: ctx.String("log-format"),
	})
}

func getCipherEngine(ctx *cli.Context) (cipher.Engine, error) {
	cmds := ctx.Command.Names()
	switch cmds[len(cmds)-1] {
	case "aes":
		return cipher.NewAESClient(ctx.String("key"))
	case "aws":
		return cipher.NewAWSClient(ctx.String("kms-key-arn"))
	case "gcp":
		return cipher.NewGCPClient(ctx.String("kms-key-name"))
	case "pgp":
		return cipher.NewPGPClient(ctx.String("public-key-path"), ctx.String("private-key-path"))
	case "vault":
		return cipher.NewVaultClient(ctx.String("transit-key"))
	default:
		return nil, fmt.Errorf("Engine %v is not implemented yet", ctx.Command.FullName())
	}
}

func readInput(ctx *cli.Context) (input string, err error) {
	switch ctx.NArg() {
	case 0:
		read, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		input = string(read)
	case 1:
		input = ctx.Args().First()
	default:
		return "", fmt.Errorf("Too many arguments provided")
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

	return cli.NewExitError("", exitCode)
}

// ExecWrapper gracefully logs and exits our `run` functions
func ExecWrapper(f func(ctx *cli.Context) (int, error)) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		if err := mlock.LockMemory(); err != nil {
			return exit(1, fmt.Errorf("error locking s5 memory: %w", err))
		}
		return exit(f(ctx))
	}
}
