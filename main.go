package main

import (
	"github.com/andreaspenz/shadow/cli"
	"os"
)

var (
	// version is overridden at linking time
	version = "dev"
	// channel is overridden at linking time
	channel = "dev"
	// overridden at linking time
	buildDate string
)

func main() {
	if err := cli.NewApplication(version, channel, buildDate).Run(os.Args); err != nil {
		os.Exit(-1)
	}
}
