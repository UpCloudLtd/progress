package messages

import "fmt"

type Color interface {
	Sprint(...interface{}) string
	Sprintf(string, ...interface{}) string
}

type noColor struct{}

func (noColor) Sprint(args ...interface{}) string {
	return fmt.Sprint(args...)
}

func (noColor) Sprintf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
