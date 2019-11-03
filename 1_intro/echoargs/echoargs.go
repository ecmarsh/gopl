// Prints command-line arguments to standard output
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	n = flag.Bool("n", false, "omit trailing newline")
	s = flag.String("s", " ", "separator")
)

// By having echoargs write through this variable, and
// not directly to `os.Stdout`, the tests can substitute
// a different Writer implementation that records what
// was written for later inspection.
var out io.Writer = os.Stdout

// main parses and reads flag values, then reports errors from echo.
func main() {
	flag.Parse()
	if err := echoargs(!*n, *s, flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "echo %v\n", err)
		os.Exit(1)
	}
}

// echoargs
func echoargs(newline bool, sep string, args []string) error {
	fmt.Fprintf(out, strings.Join(args, sep))
	if newline {
		fmt.Fprintln(out)
	}
	return nil
}

/*
// echo args v1
func main() {
	var s, sep string
	for i, arg := range os.Args[1:] {
		n := strconv.FormatInt(int64(i+1), 10)
		s += sep + "(" + n + ")" + arg
		sep = " "
	}
	fmt.Println(s)
}
*/
