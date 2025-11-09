package main

type RedirectType int

const (
	RedirectStdout RedirectType = iota
	RedirectStderr
)

func findRedirectOutputIndex(commands []string) (int, bool, RedirectType) {
	// - allow to find redirect 'flag' index
	for i, p := range commands {
		switch p {
		case ">", "1>":
			redirectIndex := i
			return redirectIndex, true, RedirectStdout
		case "2>":
			redirectIndex := i
			return redirectIndex, true, RedirectStderr
		}
	}
	return -1, false, RedirectStdout
}
