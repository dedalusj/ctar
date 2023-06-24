//go:build windows

package main

import (
	"os"
	"syscall"
	"time"
)

func fileInfoToFile(path string, info os.FileInfo) File {
	mtime := info.ModTime()
	fileTime := info.Sys().(*syscall.Win32FileAttributeData).LastAccessTime
	atime := time.Unix(0, fileTime.Nanoseconds())

	return File{
		Path:  path,
		Size:  info.Size(),
		Atime: atime,
		Mtime: mtime,
		Info:  info,
	}
}
