package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Stdout, os.Stderr, os.Args))
}

func run(stdOut io.Writer, stdErr io.Writer, args []string) int {
	parsedArgs, exitCode := parseCliArgs(stdErr, args)
	if exitCode != nil {
		return *exitCode
	}

	if err := Archive(parsedArgs, stdOut); err != nil {
		_, _ = fmt.Fprintf(stdErr, "Failed to create the archive: %s", err.Error())
		return 1
	}

	return 0
}
