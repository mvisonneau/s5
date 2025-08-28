package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/vault/sdk/helper/mlock"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"

	"github.com/mvisonneau/s5/internal/logs"
	"github.com/mvisonneau/s5/pkg/cipher"
)

func configure(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	var (
		logger   *zap.Logger
		encoding string
		err      error
	)

	logger, encoding, err = logs.NewLogger(cmd.String("log-level"), cmd.String("log-format"))
	if err != nil {
		return ctx, err
	}

	ctx = logs.StoreLoggerInContext(ctx, logger, encoding)

	return ctx, nil
}

func getCipherEngine(ctx context.Context, cmd *cli.Command) (engine cipher.Engine, err error) {
	cmds := cmd.Names()
	switch cmds[len(cmds)-1] {
	case "aes":
		engine, err = cipher.NewAESClient(cmd.String("key"))
	case "aws":
		engine, err = cipher.NewAWSClient(ctx, cmd.String("kms-key-arn"))
	case "gcp":
		engine, err = cipher.NewGCPClient(ctx, cmd.String("kms-key-name"))
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
		var err error

		if ctx, err = configure(ctx, cmd); err != nil {
			logs.LoggerFromContext(ctx).Fatal("failed to configure", zap.Error(err))
		}

		if err = mlock.LockMemory(); err != nil {
			logs.LoggerFromContext(ctx).Fatal("s5 requires the IPC_LOCK capability in order to secure its memory", zap.Error(err))
		}

		return f(ctx, cmd)
	}
}
