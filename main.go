package main

import (
	"os"

	"github.com/romboooo/ttracker/cli"
)

func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
