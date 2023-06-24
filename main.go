package main

import (
	"fmt"
	"os"
)

func main() {
	cliArgs := parseCliArgs()

	if err := Archive(cliArgs); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create the archive: %s", err.Error())
		os.Exit(1)
	}
}
