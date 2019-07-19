package command

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/mvisonneau/s5/logger"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var start time.Time

// Command is an interface of supported/required commands for each cipher engine
type Command interface {
	configure(*cli.Context) error
	cipher(string) (string, error)
	decipher(string) (string, error)
}

func configure(cmd Command, ctx *cli.Context) error {
	start = ctx.App.Metadata["startTime"].(time.Time)

	lc := &logger.Config{
		Level:  ctx.GlobalString("log-level"),
		Format: ctx.GlobalString("log-format"),
	}

	if err := lc.Configure(); err != nil {
		return err
	}

	return cmd.configure(ctx)
}

func getCipherEngine(ctx *cli.Context) (Command, error) {
	cmds := strings.Fields(ctx.Command.FullName())
	switch cmds[len(cmds)-1] {
	case "aes":
		return &aes{}, nil
	case "aws":
		return &aws{}, nil
	case "gcp":
		return &gcp{}, nil
	case "pgp":
		return &pgp{}, nil
	case "vault":
		return &vault{}, nil
	default:
		return nil, fmt.Errorf("Engine %v is not implemented yet", ctx.Command.FullName())
	}
}

func Cipher(ctx *cli.Context) error {
	cmd, err := getCipherEngine(ctx)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	if err := configure(cmd, ctx); err != nil {
		return cli.NewExitError(err, 1)
	}

	input, err := readInput(ctx)
	if err != nil {
		cli.ShowSubcommandHelp(ctx)
		return exit(err, 1)
	}

	ciphered, err := cmd.cipher(input)
	if err != nil {
		return exit(err, 1)
	}

	fmt.Printf("{{ s5:%s }}", ciphered)

	return exit(nil, 0)
}

func Decipher(ctx *cli.Context) error {
	cmd, err := getCipherEngine(ctx)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	if err := configure(cmd, ctx); err != nil {
		return cli.NewExitError(err, 1)
	}

	input, err := readInput(ctx)
	if err != nil {
		cli.ShowSubcommandHelp(ctx)
		return exit(err, 1)
	}

	log.Debug("Validating input string")
	re := regexp.MustCompile("{{ s5:(.*) }}")
	if !re.MatchString(input) {
		return exit(fmt.Errorf("Invalid string format, should be '{{ s5:* }}'"), 1)
	}

	deciphered, err := cmd.decipher(re.FindStringSubmatch(input)[1])
	if err != nil {
		return exit(err, 1)
	}

	fmt.Print(deciphered)

	return exit(nil, 0)
}

func Render(ctx *cli.Context) error {
	cmd, err := getCipherEngine(ctx)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	if err := configure(cmd, ctx); err != nil {
		return cli.NewExitError(err, 1)
	}

	if ctx.NArg() > 1 ||
		(ctx.String("output") != "" && ctx.Bool("in-place")) {
		cli.ShowSubcommandHelp(ctx)
		return exit(fmt.Errorf("Invalid arguments"), 1)
	}

	var fi *os.File
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

	re := regexp.MustCompile("{{ s5:([A-Za-z0-9+\\/=]*) }}")
	in := bufio.NewScanner(fi)

	var buf bytes.Buffer

	log.Debug("Starting deciphering")
	for in.Scan() {
		buf.WriteString(
			re.ReplaceAllStringFunc(in.Text(), func(src string) string {
				log.Debugf("found: %v", re.FindStringSubmatch(src)[1])
				plain, err := cmd.decipher(re.FindStringSubmatch(src)[1])
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

func exit(err error, exitCode int) *cli.ExitError {
	defer log.Debugf("Executed in %s, exiting..", time.Since(start))
	if err != nil {
		log.Error(err.Error())
		return cli.NewExitError("", exitCode)
	}

	return cli.NewExitError("", 0)
}
