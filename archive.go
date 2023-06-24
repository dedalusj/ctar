package main

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"

	gzip "github.com/klauspost/pgzip"
)

type Args struct {
	SourceDir   string
	ArchiveFile string
	MaxSize     uint64
	Verbose     bool
}

func Archive(args Args, stdOut io.Writer) error {
	files, err := listFiles(args.SourceDir)
	if err != nil {
		return fmt.Errorf("failed reading files in %s: %w", args.SourceDir, err)
	}

	out, err := os.Create(args.ArchiveFile)
	if err != nil {
		return fmt.Errorf("failed creating archive file %s: %w", args.ArchiveFile, err)
	}
	defer func() { _ = out.Close() }()

	gw := gzip.NewWriter(out)
	defer func() { _ = gw.Close() }()

	tw := tar.NewWriter(gw)
	defer func() { _ = tw.Close() }()

	var totalSize uint64
	for _, file := range files {
		if args.MaxSize > 0 && uint64(file.Size)+totalSize > args.MaxSize {
			return nil
		}

		if err := addToArchive(tw, file); err != nil {
			return fmt.Errorf("failed to add file %s to archive: %w", file.Path, err)
		}

		if args.Verbose {
			_, _ = fmt.Fprintln(stdOut, file.Path)
		}

		totalSize += uint64(file.Size)
	}

	return nil
}

type File struct {
	Path  string
	Size  int64
	Atime time.Time
	Mtime time.Time
	Info  os.FileInfo
}

func listFiles(sourceDir string) ([]File, error) {
	var files []File
	err := filepath.Walk(sourceDir, func(path string, info fs.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			files = append(files, fileInfoToFile(path, info))
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Atime.After(files[j].Atime)
	})

	return files, nil
}

func addToArchive(tw *tar.Writer, file File) error {
	f, err := os.Open(file.Path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	// Create a tar Header from the FileInfo data
	header, err := tar.FileInfoHeader(file.Info, file.Info.Name())
	if err != nil {
		return err
	}

	// Use full path as name (FileInfoHeader only takes the basename)
	// If we don't do this the directory structure would not be preserved
	header.Name = file.Path

	// Write file header to the tar archive
	if err = tw.WriteHeader(header); err != nil {
		return err
	}

	// Copy file content to tar archive
	if _, err = io.Copy(tw, f); err != nil {
		return err
	}

	return nil
}
