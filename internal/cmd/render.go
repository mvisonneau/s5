package cmd

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"regexp"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"

	"github.com/mvisonneau/s5/internal/logs"
	"github.com/mvisonneau/s5/pkg/cipher"
)

// Render is used for the render commands.
func Render(ctx context.Context, cmd *cli.Command) error {
	logger := logs.LoggerFromContext(ctx)
	cipherEngine, err := getCipherEngine(ctx, cmd)
	if err != nil {
		return err
	}

	if cmd.NArg() > 1 ||
		(cmd.String("output") != "" && cmd.Bool("in-place")) {
		if err = cli.ShowSubcommandHelp(cmd); err != nil {
			return errors.Wrap(err, "rendering subcommand help")
		}

		return errors.New("invalid arguments")
	}

	var fi *os.File

	if cmd.NArg() == 1 {
		logger.Debug("opening input file", zap.String("file", cmd.Args().First()))

		fi, err = os.Open(cmd.Args().First())
		if err != nil {
			return errors.Wrap(err, "opening input file")
		}
	} else {
		logger.Debug("reading from stdin")

		fi = os.Stdin
	}

	re := regexp.MustCompile(cipher.InputRegexp)
	in := bufio.NewScanner(fi)

	var buf bytes.Buffer

	logger.Debug("starting deciphering")

	for in.Scan() {
		buf.WriteString(
			re.ReplaceAllStringFunc(in.Text(), func(src string) string {
				logger.Debug("found content to decipher", zap.String("content", re.FindStringSubmatch(src)[1]))

				plain, err := cipherEngine.Decipher(ctx, re.FindStringSubmatch(src)[1])
				if err != nil {
					panic(err)
				}

				return plain
			}) + "\n")
	}

	if err = fi.Close(); err != nil {
		return errors.Wrap(err, "closing input file")
	}

	if err = in.Err(); err != nil {
		return errors.Wrap(err, "reading input file")
	}

	var fo *os.File

	switch {
	case cmd.String("output") != "":
		logger.Debug("creating and saving to file", zap.String("file", cmd.String("output")))

		fo, err = os.Create(cmd.String("output"))
		if err != nil {
			return errors.Wrap(err, "creating output file")
		}

		defer closeFile(ctx, fo)
	case cmd.Bool("in-place"):
		logger.Debug("updating the source file (in-place)")

		fo, err = os.Create(cmd.Args().First())
		if err != nil {
			return errors.Wrap(err, "writing output file")
		}

		defer closeFile(ctx, fo)
	default:
		logger.Debug("writing to stdout")

		fo = os.Stdout
	}

	out := bufio.NewWriter(fo)

	if _, err = out.Write(buf.Bytes()); err != nil {
		return errors.Wrap(err, "writing output file")
	}

	if err = out.Flush(); err != nil {
		return errors.Wrap(err, "flushing output file")
	}

	return nil
}

func closeFile(ctx context.Context, f *os.File) {
	if err := f.Close(); err != nil {
		logs.LoggerFromContext(ctx).Error("closing file")
	}
}
