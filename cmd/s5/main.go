package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mvisonneau/s5/internal/cli"
)

func main() {
	if err := cli.NewApp().Run(context.Background(), os.Args); err != nil {
		fmt.Println(err.Error())
	}
}
