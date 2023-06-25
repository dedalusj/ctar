//go:build linux

package main

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func fileInfoToFile(path string, info os.FileInfo) File {
	mtime := info.ModTime()
	stat := info.Sys().(*syscall.Stat_t)
	atime := time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))

	return File{
		Path:  path,
		Size:  info.Size(),
		Atime: atime,
		Mtime: mtime,
		Info:  info,
	}
}

func cleanPath(p string) string {
	if !filepath.IsAbs(p) {
		return p
	}
	return strings.TrimPrefix(p, "/")
}
