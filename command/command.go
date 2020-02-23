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

// Cipher is used for the cipher commands
func Cipher(ctx *cli.Context) error {
	cipherEngine, err := getCipherEngine(ctx)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	if err := configure(ctx); err != nil {
		return cli.NewExitError(err, 1)
	}

	input, err := readInput(ctx)
	if err != nil {
		cli.ShowSubcommandHelp(ctx)
		return exit(err, 1)
	}

	if !ctx.Bool("no-trim") {
		input = strings.Trim(input, " \n")
	}

	ciphered, err := cipherEngine.Cipher(input)
	if err != nil {
		return exit(err, 1)
	}

	fmt.Print(cipher.GenerateOutput(ciphered))

	return exit(nil, 0)
}

// Decipher is used for the decipher commands
func Decipher(ctx *cli.Context) error {
	cipherEngine, err := getCipherEngine(ctx)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	if err := configure(ctx); err != nil {
		return cli.NewExitError(err, 1)
	}

	input, err := readInput(ctx)
	if err != nil {
		cli.ShowSubcommandHelp(ctx)
		return exit(err, 1)
	}

	log.Debug("Validating input string")
	parsedInput, err := cipher.ParseInput(input)
	if err != nil {
		return exit(err, 1)
	}

	deciphered, err := cipherEngine.Decipher(parsedInput)
	if err != nil {
		return exit(err, 1)
	}

	fmt.Print(deciphered)

	return exit(nil, 0)
}

// Render is used for the render commands
func Render(ctx *cli.Context) error {
	cipherEngine, err := getCipherEngine(ctx)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	if err := configure(ctx); err != nil {
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
