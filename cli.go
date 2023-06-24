package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/dustin/go-humanize"
)

// parseCliArgs parse and validate the cli args
// if any of the cli args fail validation it will print an error and exit
func parseCliArgs(stdErr io.Writer, args []string) (Args, *int) {
	programName := args[0]
	cli := flag.NewFlagSet(programName, flag.ContinueOnError)
	cli.SetOutput(stdErr)

	cli.Usage = func() {
		fmt.Fprintf(cli.Output(), "usage: %s [-sv] source-dir archive-file\n", programName)
		cli.PrintDefaults()
	}

	var maxSizeStr string
	cli.StringVar(&maxSizeStr, "s", "0", "Maximum size of all the files to include in the archive. Use 0 if all the files are to be archived")

	var verbose bool
	cli.BoolVar(&verbose, "v", false, "Verbose mode to list files included in the archive")

	if err := cli.Parse(args[1:]); err != nil {
		exitCode := ptr(1)
		if errors.Is(err, flag.ErrHelp) {
			exitCode = ptr(0)
		}
		return Args{}, exitCode
	}

	args = cli.Args()
	if len(args) < 1 {
		_, _ = fmt.Fprintln(stdErr, "source-dir missing")
		cli.Usage()
		return Args{}, ptr(1)
	}
	if len(args) < 2 {
		_, _ = fmt.Fprintln(stdErr, "archive-file missing")
		cli.Usage()
		return Args{}, ptr(1)
	}

	sourceDir, archiveFile := args[0], args[1]
	if sdInfo, err := os.Stat(sourceDir); os.IsNotExist(err) {
		_, _ = fmt.Fprintf(stdErr, "Source directory %s does not exists\n", sourceDir)
		cli.Usage()
		return Args{}, ptr(1)
	} else if !sdInfo.IsDir() {
		_, _ = fmt.Fprintf(stdErr, "Source directory %s is not a directory\n", sourceDir)
		cli.Usage()
		return Args{}, ptr(1)
	}

	maxSize, err := humanize.ParseBytes(maxSizeStr)
	if err != nil {
		_, _ = fmt.Fprintln(stdErr, "Invalid format for max size")
		cli.Usage()
		return Args{}, ptr(1)
	}

	return Args{
		SourceDir:   sourceDir,
		ArchiveFile: archiveFile,
		MaxSize:     maxSize,
		Verbose:     verbose,
	}, nil
}

func ptr[T any](v T) *T {
	return &v
}
