package command

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// Render function
func Render(ctx *cli.Context) error {
	if err := configure(ctx); err != nil {
		return cli.NewExitError(err, 1)
	}

	if ctx.NArg() > 1 ||
		(ctx.String("output") != "" && ctx.Bool("in-place")) {
		cli.ShowSubcommandHelp(ctx)
		return exit(fmt.Errorf("Invalid arguments"), 1)
	}

	var fi *os.File
	var err error
	if ctx.NArg() == 1 {
		log.Debugf("Opening input file : %s", ctx.Args().First())
		fi, err = os.Open(ctx.Args().First())
		if err != nil {
			return exit(err, 1)
		}
	} else {
		log.Debug("Reading from stdin")
		fi = os.Stdin
	}

	re := regexp.MustCompile("{{ (s5:[A-Za-z0-9+\\/=]*) }}")
	in := bufio.NewScanner(fi)

	var buf bytes.Buffer

	log.Debug("Starting deciphering")
	for in.Scan() {
		buf.WriteString(
			re.ReplaceAllStringFunc(in.Text(), func(src string) string {
				log.Debugf("found: %v", re.FindStringSubmatch(src)[1])
				plain, err := vt.Decipher(re.FindStringSubmatch(src)[1])
				if err != nil {
					panic(err)
				}
				return plain
			}) + "\n")
	}

	fi.Close()

	if err := in.Err(); err != nil {
		return exit(err, 1)
	}

	var fo *os.File
	if ctx.String("output") != "" {
		log.Debugf("Creating and outputing to file : %s", ctx.String("output"))
		fo, err = os.Create(ctx.String("output"))
		if err != nil {
			return exit(err, 1)
		}
		defer fo.Close()
	} else if ctx.Bool("in-place") {
		log.Debug("Updating the source file (in-place)")
		fo, err = os.Create(ctx.Args().First())
		if err != nil {
			return exit(err, 1)
		}
		defer fo.Close()
	} else {
		log.Debug("Outputing to stdout")
		fo = os.Stdout
	}

	out := bufio.NewWriter(fo)
	_, err = out.Write(buf.Bytes())
	if err != nil {
		return exit(err, 1)
	}

	err = out.Flush()
	if err != nil {
		return exit(err, 1)
	}

	return exit(nil, 0)
}
