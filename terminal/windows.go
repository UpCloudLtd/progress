package terminal

import (
	"io"
	"os"
	"runtime"

	"golang.org/x/term"
)

func getUnicodeSafeWindowsTermPrograms() map[string]bool {
	return map[string]bool{
		"vscode": true,
	}
}

func IsWindowsTerminal(target io.Writer) bool {
	if runtime.GOOS != "windows" {
		return false
	}

	file, ok := target.(*os.File)
	if !ok {
		return false
	}

	return term.IsTerminal(int(file.Fd())) // #nosec G115 -- File-descriptor should be safe to convert into int.
}

// IsUnicodeSafeWindowsTermProgram determines if current terminal program is likely able to output unicode characters. Terminal program is determined from TERM_PROGRAM environment variable.
func IsUnicodeSafeWindowsTermProgram() bool {
	termProg := os.Getenv("TERM_PROGRAM")
	return getUnicodeSafeWindowsTermPrograms()[termProg]
}
