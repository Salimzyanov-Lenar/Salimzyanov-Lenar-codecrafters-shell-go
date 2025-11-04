package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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

func checkPathInSystem(command string) (bool, string) {
	// Type handler for check program in system
	paths := os.Getenv("PATH")

	for _, path := range strings.Split(paths, ":") {
		file := filepath.Join(path, command)
		if _, err := os.Stat(file); err == nil {
			return true, file
		}
	}
	return false, ""
}

func checkType(commands []string) {
	// Type handler
	command := strings.TrimSpace(commands[1])
	_, ok := commandHandlers[command]
	if ok {
		fmt.Fprintln(os.Stdout, fmt.Sprint(command, " is a shell builtin"))
		return
	}
	command = strings.TrimSpace(commands[1])
	exists, filePath := checkPathInSystem(command)
	if !exists {
		fmt.Fprintln(os.Stdout, strings.TrimRight(strings.Join(commands[1:], ""), "\n")+": not found")
	}
	fmt.Fprintf(os.Stdout, "%s is %s\n", command, filePath)
}

func runExternal(path string, command string, args []string) {
	cmd := exec.Command(path, args...)
	cmd.Args = append([]string{command}, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func anotherProgramHandler(commands []string) bool {
	// Parse program name and try find it
	// 		Ex: python3 -> true, /usr/bin/python3
	command := strings.TrimSpace(commands[0])
	exists, filePath := checkPathInSystem(command)
	args := []string{}
	if len(commands) >= 2 {
		args = commands[1:]
		args[len(args)-1] = strings.TrimSpace(args[len(args)-1])
	}
	if exists {
		runExternal(filePath, command, args)
		return true
	}
	return false
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

		// Output
		commands := strings.Split(input, " ")
		handler, ok := commandHandlers[commands[0]]

		// Builtin handler
		if ok {
			handler(commands)
			// Try to run program in system
		} else if anotherProgramHandler(commands) {
			continue
			// Program not found
		} else {
			fmt.Fprintln(os.Stdout, input[:len(input)-1]+": command not found")
			continue
		}
	}
}
