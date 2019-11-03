package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestEcho(t *testing.T) {
	// test cases
	var tests = []struct {
		newline bool
		sep     string
		args    []string
		want    string // expected
	}{
		{true, "", []string{}, "\n"},
		{false, "", []string{}, ""},
		{true, "\t", []string{"one", "two", "three"}, "one\ttwo\tthree\n"},
		{true, ",", []string{"a", "b", "c"}, "a, b, c\n"},
		{false, ":", []string{"1", "2", "3"}, "1:2:3"},
	}
	// test case drive
	for _, test := range tests {
		descr := fmt.Sprintf("echoargs(%v, %q, %q)",
			test.newline, test.sep, test.args)
		out = new(bytes.Buffer) // for captured output, global in echoargs
		if err := echoargs(test.newline, test.sep, test.args); err != nil {
			t.Errorf("%s failed: %v", descr, err)
			continue
		}
		got := out.(*bytes.Buffer).String()
		if got != test.want {
			t.Errorf("%s = %q, want %q", descr, got, test.want)
		}
	}
}

/*
// Usage

$ cd $GOPATH
$ go test ./1_intro_/echoargs
--- FAIL: TestEcho (0.00s)
    echoargs_test.go:34: echoargs(true, ",", ["a" "b" "c"]) = "a,b,c\n", want "a, b, c\n"
FAIL
FAIL	_/[GOPATH]/1_intro/echoargs	0.006s
FAIL

*/
