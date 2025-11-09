package main

import (
	"fmt"
	"os"
	"strings"
)

// Package for commands
// - include builtin handlers and external
// - actually execution of every program happen in executor.go

// Case: exit <status_code> -> exit from terminal
func exitHandler(commands []string) {
	if len(commands) == 1 {
		os.Exit(0)
	}
	if len(commands) > 2 {
		fmt.Fprintln(os.Stderr, "exit: too many arguments")
		return
	}
	if len(commands) == 2 && commands[1] == "0" {
		os.Exit(0)
	}
}

// Case: echo <some_text> -> print <some_text> in terminal
func echoHandler(commands []string) {
	runBuiltIn(commands, func(cmds []string) {
		if len(cmds) >= 1 {
			fmt.Fprintln(os.Stdout, strings.Join(cmds[1:], " "))
		}
	})
}

// Case: type <program_name> -> show path or say that builtin
func typeHandler(commands []string) {
	if len(commands) < 2 {
		fmt.Fprintln(os.Stderr, "type: missing argument")
		return
	}

	cmdName := strings.TrimSpace(commands[1])
	runBuiltIn(commands, func(cmds []string) {
		if len(cmds) >= 1 {
			// find in builtin's
			if _, ok := commandHandlers[cmdName]; ok {
				fmt.Fprintf(os.Stdout, "%s is a shell builtin\n", cmdName)
				return
			}
			// find external
			exists, filePath := checkFileExistsInSystem(cmdName)
			if exists {
				fmt.Fprintf(os.Stdout, "%s is %s\n", cmdName, filePath)
			} else {
				fmt.Fprintf(os.Stdout, "%s: not found\n", cmdName)
			}
			return
		}
	})
}

// Case: pwd -> /some/path/to/current/dir
func pwdHandler(commands []string) {
	if len(commands) > 1 {
		return
	}
	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
	}
	fmt.Println(path)
}

// Case: cd <path> -> change current dir to <path>
func cdHandler(commands []string) {
	if len(commands) < 2 {
		fmt.Fprintln(os.Stderr, "cd: missing argument")
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
		fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", absolutePath)
	}
}

// Case: program not builtin, find it in PATH dirs
func externalProgramHandler(commands []string) bool {
	command := strings.TrimSpace(commands[0])
	programExists, filePath := checkFileExistsInSystem(command)
	redirectIndex, redirectExists, redirectType := findRedirectOutputIndex(commands)

	args := []string{}
	if len(commands) > 1 {
		args = commands[1:]
	}

	if programExists {
		if redirectExists {
			runExternalRedirected(filePath, command, redirectIndex, redirectType, args)
		} else {
			runExternal(filePath, command, args)
		}
		return true
	}

	return false
}

func shellestHandler(commands []string) {
	art := `
        ::                             
        -+                             
        -%%- :                         
        +%%##%+.=:                     
       =##*+==++#%=                    
       #%+:====+===: *:                
      #=*-..==-:-:-:-%%#               
      -#*=-:-+*=-:-=##*+:=-            
     +*+#*:.:-+*=--+*+=+*##=           
      #*++=::.:=*::=---:-=-+ -.        
   .   .**=.:=::+*-:..:-=-:= #+        
    =%%%**-...:-+*=:.:.:=---*%%-       
     ##*=-:::...:=*=:.::=:-#=*#=.      
   :+::=+=:-=-:..:=*-::--:-:-+=*=      
    =@%#+*=.:-++=.:++::=::.:+==#:      
     *#+=-:.:-:::-+**-+-...-+:==  :    
      %#=:.::==-:::+##*:...+-:-:+#*    
        -==-::-=+++=**#::.-+:.:-+*=    
    :-----=+:....:--+#*+::=-..---#-    
     :#%%#=-:..:.::-=*#*=:=::--:+*     
       -#%#***-...:--=#*#+=:=::=*      
        -#*=:-=**+-:::-*##=+::=*-      
          -##=::.-+****#%##*++++.      
          -#%#*+===::-=#%#@%#*=.       
             :-*#*++==+*%#%=:          
                  =%%%+-++*.           
                         +           
                          +          
                           =.        
                           -:      
                             ==.    
                               @#=...
                                   ..
	`
	fmt.Fprintln(os.Stdout, art)
}
