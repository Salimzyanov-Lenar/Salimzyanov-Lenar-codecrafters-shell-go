package main

import (
	"os"
	"path/filepath"
	"strings"
)

func checkPathInSystem(command string) (bool, string) {
	// Case: ./program /another_program
	if strings.Contains(command, "/") {
		info, err := os.Stat(command)
		if err == nil && !info.IsDir() && info.Mode().Perm()&0111 != 0 {
			return true, command
		}
		return false, ""
	}

	// Case: going through PATH
	for _, dir := range strings.Split(os.Getenv("PATH"), ":") {
		fullPath := filepath.Join(dir, command)
		info, err := os.Stat(fullPath)
		if err == nil && !info.IsDir() && info.Mode().Perm()&0111 != 0 {
			return true, fullPath
		}
	}

	return false, ""
}
