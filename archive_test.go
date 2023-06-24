package main

import (
	"archive/tar"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mholt/archiver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Archive(t *testing.T) {
	sizes := []uint64{
		1024 * 512,       // 500KB
		1024 * 1024 * 1,  // 1MB
		1024 * 1024 * 2,  // 2MB
		1024 * 1024 * 4,  // 4MB
		1024 * 1024 * 10, // 10MB
	}
	for _, maxSize := range sizes {
		maxSize := maxSize
		t.Run(fmt.Sprintf("it archives a directory up to the specified maximum size - %d", maxSize), func(t *testing.T) {
			archiveFilepath, unarchivedPath := getSingleTestPaths(t, fmt.Sprint(maxSize))

			err := Archive(Args{
				SourceDir:   sourceDir,
				ArchiveFile: archiveFilepath,
				MaxSize:     maxSize,
			})
			assert.NoError(t, err)

			err = archiver.Unarchive(archiveFilepath, unarchivedPath)
			assert.NoError(t, err)

			var numSourceFiles uint64
			assertDir(t, sourceDir, func(t *testing.T, path string, fileInfo os.FileInfo) {
				numSourceFiles += 1
			})

			var totalUntarredSize uint64
			var numArchivedFiles uint64
			assertDir(t, unarchivedPath, func(t *testing.T, path string, fileInfo os.FileInfo) {
				fixtureFilePath := strings.TrimPrefix(path, unarchivedPath)
				fixtureFileInfo, err := os.Stat(fixtureFilePath)
				assert.NoError(t, err)

				assert.Equal(t, fixtureFileInfo.Name(), fileInfo.Name())
				assert.Equal(t, fixtureFileInfo.Size(), fileInfo.Size())
				assert.Equal(t, fixtureFileInfo.Mode(), fileInfo.Mode())

				totalUntarredSize += uint64(fileInfo.Size())
				numArchivedFiles += 1
			})

			assert.True(t, totalUntarredSize <= maxSize)
			assert.True(t, numArchivedFiles <= numSourceFiles)
		})
	}

	t.Run("it archives all the files if max size is zero", func(t *testing.T) {
		archiveFilepath, unarchivedPath := getSingleTestPaths(t, "all")

		err := Archive(Args{
			SourceDir:   sourceDir,
			ArchiveFile: archiveFilepath,
			MaxSize:     0,
		})
		assert.NoError(t, err)

		err = archiver.Unarchive(archiveFilepath, unarchivedPath)
		assert.NoError(t, err)

		var numExpectedFiles uint64
		assertDir(t, sourceDir, func(t *testing.T, path string, fileInfo os.FileInfo) {
			numExpectedFiles += 1
		})

		var numArchivedFiles uint64
		assertDir(t, unarchivedPath, func(t *testing.T, path string, fileInfo os.FileInfo) {
			fixtureFilePath := strings.TrimPrefix(path, unarchivedPath)
			fixtureFileInfo, err := os.Stat(fixtureFilePath)
			assert.NoError(t, err)

			assert.Equal(t, fixtureFileInfo.Name(), fileInfo.Name())
			assert.Equal(t, fixtureFileInfo.Size(), fileInfo.Size())
			assert.Equal(t, fixtureFileInfo.Mode(), fileInfo.Mode())

			numArchivedFiles += 1
		})

		assert.Equal(t, numExpectedFiles, numArchivedFiles)
	})

	t.Run("it preserves the files modified time", func(t *testing.T) {
		archiveFilepath, _ := getSingleTestPaths(t, "mtime")

		// set the time of the source files back one week to ease the testing of mtime
		oneWeekAgo := time.Now().AddDate(0, 0, -7)
		setTimes(t, sourceDir, oneWeekAgo, oneWeekAgo)

		err := Archive(Args{
			SourceDir:   sourceDir,
			ArchiveFile: archiveFilepath,
			MaxSize:     0,
		})
		assert.NoError(t, err)

		err = archiver.Walk(archiveFilepath, func(f archiver.File) error {
			if header, ok := f.Header.(*tar.Header); ok {
				fixtureFileInfo, err := os.Stat(header.Name)
				assert.NoError(t, err)

				fixtureModTime := fixtureFileInfo.ModTime().Truncate(time.Second)
				assert.WithinDuration(t, fixtureModTime, f.ModTime(), 2*time.Second)
				assert.WithinDuration(t, fixtureModTime, header.ModTime, 2*time.Second)
			} else {
				assert.Fail(t, "Wrong header format")
			}
			return nil
		})
		assert.NoError(t, err)
	})
}

func getSingleTestPaths(t *testing.T, prefix string) (string, string) {
	archiveFilepath := path.Join(testDir, fmt.Sprintf("%s-archive.tar.gz", prefix))
	t.Cleanup(func() {
		if err := os.RemoveAll(archiveFilepath); err != nil {
			panic(err)
		}
	})

	unarchivedPath := path.Join(testDir, fmt.Sprintf("%s-unarchived", prefix))
	t.Cleanup(func() {
		if err := os.RemoveAll(unarchivedPath); err != nil {
			panic(err)
		}
	})

	return archiveFilepath, unarchivedPath
}

type fileTester func(t *testing.T, path string, fileInfo os.FileInfo)

func assertDir(t *testing.T, dirPath string, tester fileTester) {
	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		tester(t, path, info)
		return nil
	})
	assert.NoError(t, err)
}

func setTimes(t *testing.T, dirPath string, atime time.Time, mtime time.Time) {
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return os.Chtimes(path, atime, mtime)
	})
	require.NoError(t, err)
}
