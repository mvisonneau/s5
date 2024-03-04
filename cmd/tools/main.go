package main

import (
	"fmt"
	"time"

	"github.com/mvisonneau/s5/internal/cli"
)

var version = "devel"

func main() {
	fmt.Println(cli.NewApp(version, time.Now()).ToMan())
}
