package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

var commandHandlers map[string]func([]string)

// Initialize with built in function
func init() {
	commandHandlers = map[string]func([]string){
		"echo": echoHandler,
		"exit": exitHandler,
		"type": typeHandler,
		"pwd":  pwdHandler,
		"cd":   cdHandler,
	}
}

func main() {
	for {
		// Input
		fmt.Fprint(os.Stdout, "$ ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		// Parse input
		commands := parseInput(input)
		handler, ok := commandHandlers[strings.TrimSpace(commands[0])]

		// Output
		if ok {
			handler(commands) // Builtin handler
		} else if anotherProgramHandler(commands) { // Try to run program in system
			continue
		} else { // Program not found
			fmt.Fprintln(os.Stdout, input[:len(input)-1]+": command not found")
			continue
		}
	}
}

func parseInput(input string) []string {
	var current strings.Builder
	parts := []string{}
	inSingleQuote := false
	inDoubleQuote := false
	escaped := false

	// Trim trailing newline from ReadString
	input = strings.TrimSuffix(input, "\n")

	for _, c := range input {
		switch {
		case escaped:
			if inDoubleQuote {
				switch c {
				case '"', '\\':
					current.WriteRune(c)
				default:
					current.WriteRune('\\')
					current.WriteRune(c)
				}
			} else {
				current.WriteRune(c)
			}
			escaped = false

		case c == '\\' && !inDoubleQuote && !inSingleQuote:
			escaped = true
		case c == '\\' && inDoubleQuote:
			escaped = true

		case c == '\'' && !inDoubleQuote:
			inSingleQuote = !inSingleQuote
		case c == '"' && !inSingleQuote:
			inDoubleQuote = !inDoubleQuote
		case unicode.IsSpace(c) && !inSingleQuote && !inDoubleQuote:
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}

		default:
			current.WriteRune(c)
		}
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	return parts
}
