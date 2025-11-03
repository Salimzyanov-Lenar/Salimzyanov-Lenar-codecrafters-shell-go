package main

import "strings"

func handleExit(command string) bool {
	parts := strings.Fields(command)

	if len(parts) != 2 {
		return false
	}

	statusCode := parts[1]

	if parts[0] != "exit" {
		return false
	}

	if statusCode == "0" {
		return true
	}

	return false
}
