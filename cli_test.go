package main

import (
	"regexp"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseCliArgs(t *testing.T) {
	t.Run("it prints the help if requested", func(t *testing.T) {
		_, stdErr := newMockWriters()

		_, exitCode := parseCliArgs(stdErr, []string{"app", "-h"})
		assert.NotNil(t, exitCode)
		assert.Equal(t, 0, *exitCode)
		assert.Contains(t, stdErr.Content(), "Usage")
		assert.Contains(t, stdErr.Content(), "-v")
		assert.Contains(t, stdErr.Content(), "-size")
	})

	t.Run("it returns an error if both arguments are missing", func(t *testing.T) {
		_, stdErr := newMockWriters()

		_, exitCode := parseCliArgs(stdErr, []string{"app"})
		assert.NotNil(t, exitCode)
		assert.Equal(t, 1, *exitCode)
		assert.Contains(t, stdErr.Content(), "Usage")
		assert.Contains(t, stdErr.Content(), "Argument missing")
	})

	t.Run("it returns an error if the archive file is missing", func(t *testing.T) {
		_, stdErr := newMockWriters()

		_, exitCode := parseCliArgs(stdErr, []string{"app", "source"})
		assert.NotNil(t, exitCode)
		assert.Equal(t, 1, *exitCode)
		assert.Contains(t, stdErr.Content(), "Usage")
		assert.Contains(t, stdErr.Content(), "Argument missing")
	})

	t.Run("it returns an error if the source directory does not exists", func(t *testing.T) {
		_, stdErr := newMockWriters()

		_, exitCode := parseCliArgs(stdErr, []string{"app", "non_existing", "test.tar.gz"})
		assert.NotNil(t, exitCode)
		assert.Equal(t, 1, *exitCode)
		assert.Contains(t, stdErr.Content(), "Usage")
		assert.Regexp(t, regexp.MustCompile("Source directory \\S+ does not exists"), stdErr.Content())
	})

	t.Run("it returns an error if the source directory is not a directory", func(t *testing.T) {
		_, stdErr := newMockWriters()

		_, filename, _, ok := runtime.Caller(0)
		require.True(t, ok)

		_, exitCode := parseCliArgs(stdErr, []string{"app", filename, "test.tar.gz"})
		assert.NotNil(t, exitCode)
		assert.Equal(t, 1, *exitCode)
		assert.Contains(t, stdErr.Content(), "Usage")
		assert.Regexp(t, regexp.MustCompile("Source directory \\S+ is not a directory"), stdErr.Content())
	})

	t.Run("it returns an error if max size has an invalid format", func(t *testing.T) {
		_, stdErr := newMockWriters()

		_, exitCode := parseCliArgs(stdErr, []string{"app", "--size", "foo", testDir, "test.tar.gz"})
		assert.NotNil(t, exitCode)
		assert.Equal(t, 1, *exitCode)
		assert.Contains(t, stdErr.Content(), "Usage")
		assert.Contains(t, stdErr.Content(), "Invalid format for max size")
	})

	t.Run("it correctly parse the args", func(t *testing.T) {
		_, stdErr := newMockWriters()

		parsedArgs, exitCode := parseCliArgs(stdErr, []string{"app", "--size", "5.2MB", "-v", testDir, "test.tar.gz"})
		assert.Nil(t, exitCode)
		assert.Empty(t, stdErr.Content())
		assert.Equal(t, Args{
			SourceDir:   testDir,
			ArchiveFile: "test.tar.gz",
			MaxSize:     5200000,
			Verbose:     true,
		}, parsedArgs)
	})
}
