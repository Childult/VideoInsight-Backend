package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	Info.Println("hello")
	Warning.Println("hello")
	Debug.Println("hello")
	Error.Println("hello")
}
