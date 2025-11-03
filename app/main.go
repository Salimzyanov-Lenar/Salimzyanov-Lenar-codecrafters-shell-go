package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

type ShellState string

const (
	StateRead  ShellState = "read"
	StateEval  ShellState = "eval"
	StatePrint ShellState = "print"
)

var transitionOrder = []ShellState{StateRead, StateEval, StatePrint}

type Shell struct {
	state       ShellState
	commandLine string
	message     string
	shouldExit  bool
}

func NewShell() *Shell {
	return &Shell{state: StateRead, commandLine: "", message: ""}
}

func (sh *Shell) changeState() {
	idx := slices.Index(transitionOrder, sh.state)
	sh.state = transitionOrder[(idx+1)%3]
}

func (sh *Shell) read() {
	if sh.state != "read" {
		return
	}
	fmt.Printf("$ ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			sh.commandLine = ""
			sh.message = "\n"
			sh.shouldExit = true
			return
		}
	}
	line = strings.TrimRight(line, "\r\n")
	sh.commandLine = line
}

func isExecutable(mode os.FileMode) bool {
	/*
		0111 => 0001001001
		rwx, so 0th bit is x, hence this AND captures all x
		all x because executable by whom not specified
	*/
	return mode&0111 != 0
}

func (sh *Shell) eval() {
	if sh.state != StateEval {
		sh.message = "internal shell error"
		return
	}
	args := strings.Split(sh.commandLine, " ")
	switch args[0] {
	case "exit":
		{
			if len(args) < 2 {
				sh.message = "invalid arguments to exit\n"
				return
			}
			status, err := strconv.Atoi(args[1])
			if err != nil {
				sh.message = "invalid arguments to exit\n"
			} else {
				os.Exit(status)
			}
		}
	case "echo":
		{
			fmt.Printf("%s\n", strings.Join(args[1:], " "))
		}
	case "type":
		{
			if len(args) != 2 {
				sh.message = fmt.Sprintf("type requires 1 argument, %d given\n", len(args)-1)
				return
			}
			builtinCommands := map[string]struct{}{
				"exit": {},
				"echo": {},
				"type": {},
			}
			_, exists := builtinCommands[args[1]]
			if exists {
				sh.message = fmt.Sprintf("%s is a shell builtin\n", args[1])
			} else {
				sh.message = fmt.Sprintf("%s: not found\n", args[1])
				envPath := os.Getenv("PATH")
				dirs := filepath.SplitList(envPath)
				for _, dir := range dirs {
					_, err := os.Stat(dir)
					if err != nil {
						continue
					}
					// fmt.Println(dir)
					filePath := filepath.Join(dir, args[1])
					info, err := os.Stat(filePath)
					if err != nil {
						continue
					}
					if !isExecutable(info.Mode()) {
						continue
					}
					sh.message = fmt.Sprintf("%s is %s\n", args[1], filePath)
					break
				}
			}
		}
	default:
		{
			sh.message = fmt.Sprintf("%s: command not found\n", args[0])
		}
	}
}

func (sh *Shell) print() {
	if sh.state != StatePrint {
		fmt.Printf("internal shell error\n")
	}
	fmt.Printf("%s", sh.message)
}

func (sh *Shell) work() {
	switch sh.state {
	case StateRead:
		{
			sh.read()
		}
	case StateEval:
		{
			sh.eval()
		}
	case StatePrint:
		{
			sh.print()
		}
	}
}

func (sh *Shell) repl() {
	transitionCnt := len(transitionOrder)
	for range transitionCnt {
		sh.work()
		sh.changeState()
	}
}

func main() {
	sh := NewShell()
	for !sh.shouldExit {
		sh.repl()
	}
}
