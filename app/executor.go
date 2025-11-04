package main

import (
	"fmt"
	"os"
	"os/exec"
)

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
		fmt.Fprintln(os.Stderr, err)
	}
}
