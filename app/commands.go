package main

import (
	"fmt"
	"os"
	"strings"
)

func exitHandler(commands []string) {
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

func echoHandler(commands []string) {
	// Case:
	// $ echo hello world
	// hello world
	if len(commands) >= 1 {
		fmt.Fprint(os.Stdout, strings.Join(commands[1:], " "))
		return
	}
}

func typeHandler(commands []string) {
	// Case:
	// $ type echo
	// echo is a shell builtin
	// $ type ls
	// ls is /some/path/ls
	if len(commands) < 2 {
		fmt.Fprintln(os.Stderr, "type: missing argument")
		return
	}

	cmdName := strings.TrimSpace(commands[1])

	// Case: command is builtin
	if _, ok := commandHandlers[cmdName]; ok {
		fmt.Fprintf(os.Stdout, "%s is a shell builtin\n", cmdName)
		return
	}

	// Case: try to find command in PATH dirs
	exists, filePath := checkFileExistsInSystem(cmdName)
	if exists {
		fmt.Fprintf(os.Stdout, "%s is %s\n", cmdName, filePath)
	} else {
		fmt.Fprintf(os.Stdout, "%s: not found\n", cmdName)
	}
}

func pwdHandler(commands []string) {
	// Case: pwd -> /some/path/to/current/dir
	if len(commands) > 1 && commands[0] != "pwd" {
		return
	}
	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}
	fmt.Println(path)
}

func cdHandler(commands []string) {
	cmd := strings.TrimSpace(commands[0])
	if cmd != "cd" {
		return
	}

	if len(commands) < 2 {
		fmt.Println("cd: missing argument")
		return
	}

	if strings.TrimSpace(commands[1]) == "~" {
		homePath := os.Getenv("HOME")
		os.Chdir(homePath)
		return
	}

	absolutePath := strings.TrimSpace(commands[1])
	err := os.Chdir(absolutePath)
	if err != nil {
		fmt.Fprintf(os.Stdout, "cd: %s: No such file or directory\n", absolutePath)
	}
}

func anotherProgramHandler(commands []string) bool {
	// Case: when program isn't builtin
	command := strings.TrimSpace(commands[0])
	exists, filePath := checkFileExistsInSystem(command)
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
