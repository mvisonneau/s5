package main

import (
	"os"

	"github.com/mvisonneau/s5/internal/cli"
)

var version = ""

func main() {
	cli.Run(version, os.Args)
}
