package main

import (
	"log"
	"os"
	"time"

	"github.com/mvisonneau/s5/cli"
)

var version = ""

func main() {
	if err := cli.Init(&version, time.Now()).Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
