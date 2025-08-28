package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/vault/sdk/helper/mlock"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/mvisonneau/go-helpers/logger"
	"github.com/mvisonneau/s5/pkg/cipher"
)

func configure(cmd *cli.Command) error {
	if err := logger.Configure(logger.Config{
		Level:  cmd.String("log-level"),
		Format: cmd.String("log-format"),
	}); err != nil {
		return errors.Wrap(err, "configuring logger")
	}

	return nil
}

func getCipherEngine(cmd *cli.Command) (engine cipher.Engine, err error) {
	cmds := cmd.Names()
	switch cmds[len(cmds)-1] {
	case "aes":
		engine, err = cipher.NewAESClient(cmd.String("key"))
	case "aws":
		engine, err = cipher.NewAWSClient(cmd.String("kms-key-arn"))
	case "gcp":
		engine, err = cipher.NewGCPClient(cmd.String("kms-key-name"))
	case "pgp":
		engine, err = cipher.NewPGPClient(cmd.String("public-key-path"), cmd.String("private-key-path"))
	case "vault":
		engine, err = cipher.NewVaultClient(cmd.String("transit-key"))
	default:
		err = fmt.Errorf("engine %v is not implemented yet", cmd.FullName())
	}

	if err != nil {
		return nil, errors.Wrap(err, "getting cipher engine")
	}

	return
}

func readInput(cmd *cli.Command) (input string, err error) {
	switch cmd.NArg() {
	case 0:
		read, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", errors.Wrap(err, "reading from stdin")
		}

		input = string(read)
	case 1:
		input = cmd.Args().First()
	default:
		return "", errors.New("too many arguments provided")
	}

	return
}

// Execute gracefully logs and exits our `run` functions.
func Execute(f func(ctx context.Context, cmd *cli.Command) error) func(context.Context, *cli.Command) error {
	return func(ctx context.Context, cmd *cli.Command) error {
		if err := mlock.LockMemory(); err != nil {
			log.WithError(err).Warn("s5 requires the IPC_LOCK capability in order to secure its memory")
		}

		return f(ctx, cmd)
	}
}
