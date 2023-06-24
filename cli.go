package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dustin/go-humanize"
)

// parseCliArgs parse and validate the cli args
// if any of the cli args fail validation it will print an error and exit
func parseCliArgs() Args {
	var maxSizeStr string
	flag.StringVar(&maxSizeStr, "size", "0", "Maximum size of all the files to include in the archive. Use 0 if all the files are to be archived")

	var verbose bool
	flag.BoolVar(&verbose, "v", false, "List files included in the archive")

	var help bool
	flag.BoolVar(&help, "help", false, "Show usage")

	flag.Parse()

	if help {
		printUsage()
		os.Exit(0)
	}

	if len(os.Args) < 3 {
		_, _ = fmt.Fprintln(os.Stderr, "Argument missing")
		printUsage()
		os.Exit(1)
	}

	sourceDir, archiveFile := os.Args[1], os.Args[2]
	if sdInfo, err := os.Stat(sourceDir); os.IsNotExist(err) {
		_, _ = fmt.Fprintf(os.Stderr, "Source directory %s does not exists\n", sourceDir)
		printUsage()
		os.Exit(1)
	} else if !sdInfo.IsDir() {
		_, _ = fmt.Fprintf(os.Stderr, "Source directory %s is not a directory\n", sourceDir)
		printUsage()
		os.Exit(1)
	}

	maxSize, err := humanize.ParseBytes(maxSizeStr)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Invalid format for max size")
		printUsage()
		os.Exit(1)
	}

	return Args{
		SourceDir:   sourceDir,
		ArchiveFile: archiveFile,
		MaxSize:     maxSize,
		Verbose:     verbose,
	}
}

func printUsage() {
	_, _ = fmt.Fprintln(os.Stderr, "Usage")
}
