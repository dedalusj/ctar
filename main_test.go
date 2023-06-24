package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/mholt/archiver/v3"
	"github.com/stretchr/testify/assert"
)

var testDir string
var sourceDir string

func TestMain(m *testing.M) {
	packageDir := getPackageDir()
	setup(packageDir)
	exitVal := m.Run()
	teardown()
	os.Exit(exitVal)
}

func getPackageDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	return path.Dir(filename)
}

func getTempDir(packageDir string) string {
	return path.Join(packageDir, ".test", fmt.Sprintf("%s-%d", time.Now().Format("20060102"), rand.Intn(1000)))
}

func setup(packageDir string) {
	testArchiveFile := path.Join(packageDir, "fixtures", "test.tar.gz")

	testDir = getTempDir(packageDir)
	sourceDir = path.Join(testDir, "source")

	err := archiver.Unarchive(testArchiveFile, sourceDir)
	if err != nil {
		panic(err)
	}
}

func teardown() {
	if err := os.RemoveAll(testDir); err != nil {
		panic(err)
	}
}

type mockWriter struct {
	bytes.Buffer
}

func (m *mockWriter) Content() string {
	return strings.Trim(m.String(), "\n ")
}

func (m *mockWriter) Lines() []string {
	return strings.Split(m.Content(), "\n")
}

func newMockWriters() (*mockWriter, *mockWriter) {
	return &mockWriter{}, &mockWriter{}
}

func Test_run(t *testing.T) {
	t.Run("it runs the tool", func(t *testing.T) {
		stdOut, stdErr := newMockWriters()

		archiveFilepath, unarchivedPath := getSingleTestPaths(t, "run")

		args := []string{"ctar", archiveFilepath, sourceDir}

		exitCode := run(stdOut, stdErr, args)
		assert.Equal(t, 0, exitCode)

		err := archiver.Unarchive(archiveFilepath, unarchivedPath)
		assert.NoError(t, err)

		var numExpectedFiles uint64
		assertDir(t, sourceDir, func(t *testing.T, path string, fileInfo os.FileInfo) {
			numExpectedFiles += 1
		})

		var numArchivedFiles uint64
		assertDir(t, unarchivedPath, func(t *testing.T, path string, fileInfo os.FileInfo) {
			numArchivedFiles += 1
		})

		assert.Equal(t, numExpectedFiles, numArchivedFiles)
	})

	t.Run("it exits with non zero if the cli args are wrong", func(t *testing.T) {
		stdOut, stdErr := newMockWriters()

		args := []string{"ctar"}

		exitCode := run(stdOut, stdErr, args)
		assert.Equal(t, 1, exitCode)
	})

	t.Run("it exits with non zero if failing to archive", func(t *testing.T) {
		stdOut, stdErr := newMockWriters()

		archiveFilepath, _ := getSingleTestPaths(t, "dash/prefix")

		args := []string{"ctar", archiveFilepath, sourceDir}

		exitCode := run(stdOut, stdErr, args)
		assert.Equal(t, 1, exitCode)
	})
}
