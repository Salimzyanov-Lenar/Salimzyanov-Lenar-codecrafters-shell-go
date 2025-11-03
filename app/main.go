package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var _ = os.Stdout

func handleExit(command string) bool {
	parts := strings.Fields(command)

	if len(parts) != 2 {
		return false
	}

	cmd := parts[0]
	statusCode := parts[1]

	if cmd != "exit" {
		return false
	}

	if statusCode == "0" {
		return true
	}

	return false

}

func main() {
	// REPL
	for {
		fmt.Fprint(os.Stdout, "$ ")
		command, _ := bufio.NewReader(os.Stdin).ReadString('\n')

		if handleExit(command) {
			break
		}

		command = strings.TrimSpace(command)
		fmt.Printf("%s: command not found\n", command)
	}
}
