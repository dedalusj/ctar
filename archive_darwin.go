//go:build darwin

package main

import (
	"os"
	"syscall"
	"time"
)

func fileInfoToFile(path string, info os.FileInfo) File {
	mtime := info.ModTime()
	stat := info.Sys().(*syscall.Stat_t)
	atime := time.Unix(stat.Atimespec.Sec, stat.Atimespec.Nsec)

	return File{
		Path:  path,
		Size:  info.Size(),
		Atime: atime,
		Mtime: mtime,
		Info:  info,
	}
}
