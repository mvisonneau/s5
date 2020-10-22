package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"

	"github.com/mvisonneau/s5/pkg/cipher"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Render is used for the render commands
func Render(ctx *cli.Context) (int, error) {
	cipherEngine, err := getCipherEngine(ctx)
	if err != nil {
		return 1, err
	}

	if err := configure(ctx); err != nil {
		return 1, err
	}

	if ctx.NArg() > 1 ||
		(ctx.String("output") != "" && ctx.Bool("in-place")) {
		if err = cli.ShowSubcommandHelp(ctx); err != nil {
			return 1, err
		}
		return 1, fmt.Errorf("Invalid arguments")
	}

	var fi *os.File
	if ctx.NArg() == 1 {
		log.Debugf("Opening input file : %s", ctx.Args().First())
		fi, err = os.Open(ctx.Args().First())
		if err != nil {
			return 1, err
		}
	} else {
		log.Debug("Reading from stdin")
		fi = os.Stdin
	}

	re := regexp.MustCompile(cipher.InputRegexp)
	in := bufio.NewScanner(fi)

	var buf bytes.Buffer

	log.Debug("Starting deciphering")
	for in.Scan() {
		buf.WriteString(
			re.ReplaceAllStringFunc(in.Text(), func(src string) string {
				log.Debugf("found: %v", re.FindStringSubmatch(src)[1])
				plain, err := cipherEngine.Decipher(re.FindStringSubmatch(src)[1])
				if err != nil {
					panic(err)
				}
				return plain
			}) + "\n")
	}

	if err = fi.Close(); err != nil {
		return 1, err
	}

	if err = in.Err(); err != nil {
		return 1, err
	}

	var fo *os.File
	if ctx.String("output") != "" {
		log.Debugf("Creating and outputing to file : %s", ctx.String("output"))
		fo, err = os.Create(ctx.String("output"))
		if err != nil {
			return 1, err
		}
		defer closeFile(fo)
	} else if ctx.Bool("in-place") {
		log.Debug("Updating the source file (in-place)")
		fo, err = os.Create(ctx.Args().First())
		if err != nil {
			return 1, err
		}
		defer closeFile(fo)
	} else {
		log.Debug("Outputing to stdout")
		fo = os.Stdout
	}

	out := bufio.NewWriter(fo)
	_, err = out.Write(buf.Bytes())
	if err != nil {
		return 1, err
	}

	err = out.Flush()
	if err != nil {
		return 1, err
	}

	return 0, nil
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Error("closing file")
	}
}
