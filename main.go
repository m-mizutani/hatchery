package main

import (
	"os"

	"github.com/m-mizutani/hatchery/pkg/controller/cli"
)

func main() {
	if err := cli.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
