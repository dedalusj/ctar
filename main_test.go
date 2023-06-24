package main

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/mholt/archiver/v3"
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
