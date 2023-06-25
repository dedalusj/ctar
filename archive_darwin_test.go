//go:build darwin

package main

import "path/filepath"

func restoreAbsPath(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return "/" + p
}
