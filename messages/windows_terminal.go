package messages

import "os"

var unicodeSafeWindowsTermPrograms = map[string]bool{
	"vscode": true,
}

func isUnicodeSafeWindowsTermProgram() bool {
	termProg := os.Getenv("TERM_PROGRAM")
	return unicodeSafeWindowsTermPrograms[termProg]
}
