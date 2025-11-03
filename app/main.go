package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var commandHandlers map[string]func([]string)

func init() {
	commandHandlers = map[string]func([]string){
		"echo": echo,
		"exit": exit,
		"type": checkType,
	}
}

func exit(commands []string) {
	// Exit handler
	if len(commands) == 1 {
		os.Exit(0)
	}
	if len(commands) > 2 {
		fmt.Fprintln(os.Stdout, "exit: too many arguments")
		return
	}
	if len(commands) == 2 && commands[1] == "0\n" {
		os.Exit(0)
	}
}

func echo(commands []string) {
	// Echo handler
	if len(commands) >= 1 {
		fmt.Fprint(os.Stdout, strings.Join(commands[1:], " "))
		return
	}
}

func checkPathInSystem(commands []string) bool {
	// Type handler for check program in system
	paths := strings.Split(os.Getenv("PATH"), ":")

	for _, path := range paths {
		filePath := filepath.Join(path, strings.TrimSpace(commands[1]))
		if _, err := os.Stat(filePath); err == nil {
			mode := filePath.Mode()
			if mode&0111 != 0 {
				continue
			}
			command := strings.TrimSuffix(commands[1], "\n")
			fmt.Fprintf(os.Stdout, "%s is %s\n", command, filePath)
			return true
		}
	}
	return false
}

func checkType(commands []string) {
	// Type handler
	command := strings.TrimSpace(commands[1])
	_, ok := commandHandlers[command]
	if ok {
		fmt.Fprintln(os.Stdout, fmt.Sprint(command, " is a shell builtin"))
		return
	}
	if !checkPathInSystem(commands) {
		fmt.Fprintln(os.Stdout, strings.TrimRight(strings.Join(commands[1:], ""), "\n")+": not found")
	}
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		commands := strings.Split(input, " ")
		handler, ok := commandHandlers[commands[0]]
		if !ok {
			fmt.Fprintln(os.Stdout, input[:len(input)-1]+": command not found")
			continue
		}
		handler(commands)
	}
}
