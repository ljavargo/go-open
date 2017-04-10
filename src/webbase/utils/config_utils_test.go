package utils

import "testing"

func TestPrintConfig(t *testing.T) {
	PrintConfig(&struct {
		A int
		B string
		C bool
	}{A: 1, B: "test", C: true})
}
