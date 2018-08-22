package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var start time.Time
var v Vault

func run(c *cli.Context) error {
	start = time.Now()
	configureLogging(cfg.Log.Level, cfg.Log.Format)

	log.Debug("Configuring Vault")
	err := v.Configure(cfg.Vault.Address, cfg.Vault.Token, cfg.TransitKey)
	if err != nil {
		return exit(cli.NewExitError(err.Error(), 1))
	}

	log.Debugf("Executing function '%s'", c.Command.FullName())
	switch c.Command.FullName() {
	case "cipher":
		err = cipher(c)
	case "decipher":
		err = decipher(c)
	case "render":
		err = render(c)
	default:
		log.Fatalf("Function %v is not implemented", c.Command.FullName())
	}

	return exit(err)
}

func readInput(c *cli.Context) (input string, err error) {
	switch c.NArg() {
	case 0:
		read, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		input = string(read)
	case 1:
		input = c.Args().First()
	default:
		return "", fmt.Errorf("Too many arguments provided")
	}
	return
}

func cipher(c *cli.Context) error {
	input, err := readInput(c)
	if err != nil {
		cli.ShowSubcommandHelp(c)
		return exit(cli.NewExitError(err.Error(), 1))
	}

	log.Debug("Ciphering using Vault transit key")
	ciphered, err := v.Cipher(input)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	fmt.Printf("{{ %s }}", ciphered)

	return nil
}

func decipher(c *cli.Context) error {
	input, err := readInput(c)
	if err != nil {
		cli.ShowSubcommandHelp(c)
		return exit(cli.NewExitError(err.Error(), 1))
	}

	log.Debug("Validating input string")
	re := regexp.MustCompile("{{ (s5:.*) }}")
	if !re.MatchString(input) {
		return cli.NewExitError("Invalid string format, should be '{{ s5:* }}'", 1)
	}

	log.Debugf("Deciphering '%s' using Vault transit key", re.FindStringSubmatch(input)[1])
	plain, err := v.Decipher(re.FindStringSubmatch(input)[1])
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	fmt.Print(plain)

	return nil
}

func render(c *cli.Context) error {
	if c.NArg() > 1 ||
		(c.String("output") != "" && c.Bool("in-place")) {
		cli.ShowSubcommandHelp(c)
		return cli.NewExitError("", 1)
	}

	var fi *os.File
	var err error
	if c.NArg() == 1 {
		log.Debugf("Opening input file : %s", c.Args().First())
		fi, err = os.Open(c.Args().First())
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
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
				plain, err := v.Decipher(re.FindStringSubmatch(src)[1])
				if err != nil {
					panic(err)
				}
				return plain
			}) + "\n")
	}

	fi.Close()

	if err := in.Err(); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	var fo *os.File
	if c.String("output") != "" {
		log.Debugf("Creating and outputing to file : %s", c.String("output"))
		fo, err = os.Create(c.String("output"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		defer fo.Close()
	} else if c.Bool("in-place") {
		log.Debug("Updating the source file (in-place)")
		fo, err = os.Create(c.Args().First())
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		defer fo.Close()
	} else {
		log.Debug("Outputing to stdout")
		fo = os.Stdout
	}

	out := bufio.NewWriter(fo)
	_, err = out.Write(buf.Bytes())
	if err != nil {
		return err
	}

	err = out.Flush()
	if err != nil {
		return err
	}

	return nil
}

func exit(err error) error {
	log.Debugf("Executed in %s, exiting..", time.Since(start))
	return err
}
