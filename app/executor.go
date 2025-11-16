package main

import (
	"fmt"
	"os"
	"os/exec"
)

// Package for running programs
//
// - allow to run external programs
//   - allow to run external with redirect stdout
//
// - allow to run builtin programs
//   - allow to run builtit programs with redirect stdout

func runExternal(path string, command string, args []string) {
	// Run another program inside terminal
	// Args:
	//		path : path to program
	//		command : executable file name, like python3 or ls
	//		args : array of arguments for executable command
	cmd := exec.Command(path, args...)
	cmd.Args = append([]string{command}, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
	}
}

func runExternalRedirected(path string, command string, redirectIndex int, redirectType RedirectType, args []string) {
	// For not builtin commands
	if redirectIndex >= len(args) {
		fmt.Fprintln(os.Stderr, "syntax error: missing file after redirection")
		return
	}

	filename := args[redirectIndex]

	file, err := os.Create(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating file:", err)
		return
	}
	defer file.Close()

	cmdArgs := args[:redirectIndex-1] // everything before `>`
	cmd := exec.Command(path, cmdArgs...)
	cmd.Args = append([]string{command}, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if redirectType == RedirectStdout {
		cmd.Stdout = file
		cmd.Stderr = os.Stderr
	}
	if redirectType == RedirectStderr {
		cmd.Stdout = os.Stdout
		cmd.Stderr = file
	}

	err = cmd.Run()
	if err != nil {

	}
}

func runBuiltIn(commands []string, handler func([]string)) {
	// For builtin commands
	// - if have redirect flag - will redirect output
	// - if don't have redirect flag - will write to default stdout
	redirectIndex, exists, _ := findRedirectOutputIndex(commands)
	if !exists {
		handler(commands)
		return
	}

	if redirectIndex+1 > len(commands) {
		fmt.Fprintln(os.Stderr, "syntax error: missing file after '>'")
		return
	}

	filename := commands[redirectIndex+1]
	file, err := os.Create(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating file", err)
		return
	}
	defer file.Close()

	oldStdout := os.Stdout
	os.Stdout = file
	defer func() { os.Stdout = oldStdout }()

	handler(commands[:redirectIndex])

}
