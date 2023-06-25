//go:build windows

package main

import "path/filepath"

func restoreAbsPath(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	path, _ := os.Getwd()
	volumeName := filepath.VolumeName(path)
	return volumeName + "\\" + p
}
