package main

import (
	"flag"
	"os"

	"github.com/nicoaudy/cogmit/cmd"
)

func main() {
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showVersion, "v", false, "Show version information")
	flag.Parse()

	if showVersion {
		printVersion()
		os.Exit(0)
	}

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
