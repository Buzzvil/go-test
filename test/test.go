package test

import (
	"fmt"
	"runtime/debug"
	"testing"
)

func AssertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	message += fmt.Sprintf("\n%v != %v\n%s", a, b, debug.Stack())
	t.Fatal(message)
}
