package terminal

import (
	"io"
	"os"
	"runtime"

	"golang.org/x/term"
)

var unicodeSafeWindowsTermPrograms = map[string]bool{
	"vscode": true,
}

func IsWindowsTerminal(target io.Writer) bool {
	if runtime.GOOS != "windows" {
		return false
	}

	file, ok := target.(*os.File)
	if !ok {
		return false
	}

	return term.IsTerminal(int(file.Fd()))
}

func IsSafeWindowsTermProgram() bool {
	termProg := os.Getenv("TERM_PROGRAM")
	return unicodeSafeWindowsTermPrograms[termProg]
}
