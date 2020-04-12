package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/mvisonneau/go-helpers/logger"
	"github.com/mvisonneau/s5/cipher"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var start time.Time

func configure(ctx *cli.Context) error {
	start = ctx.App.Metadata["startTime"].(time.Time)

	lc := &logger.Config{
		Level:  ctx.GlobalString("log-level"),
		Format: ctx.GlobalString("log-format"),
	}

	return lc.Configure()
}

func getCipherEngine(ctx *cli.Context) (cipher.Engine, error) {
	cmds := strings.Fields(ctx.Command.FullName())
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

func exit(exitCode int, err error) *cli.ExitError {
	defer log.Debugf("Executed in %s, exiting..", time.Since(start))
	if err != nil {
		log.Error(err.Error())
	}

	return cli.NewExitError("", exitCode)
}

// ExecWrapper gracefully logs and exits our `run` functions
func ExecWrapper(f func(ctx *cli.Context) (int, error)) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		return exit(f(ctx))
	}
}
