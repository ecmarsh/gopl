// XMLselect prints the text of selected elements of an XML document.
package main

import (
	"encoding/xml" // Check out this package for more
	"fmt"
	"io"
	"os"
	"strings"
)

// xmlselect extracts and prints the text found beneath
// certain elements in an XML document tree.
// Each time loop encounters a `Start Element`, it pushes
// the element's name onto a stack, and for each EndElement
// it pops the name from the stack.
// API guarantees that the sequence of startElement and EndElement
// tokens will be matchd properly, even in ill-formed documents.
// Comments are ignored and when xmlselect encounters a CharData,
// it prints the text only if the stack contains all the elements
// named by the command-line arguments, in order.
func main() {
	dec := xml.NewDecoder(os.Stdin)
	var stack []string // stack of element names
	for {
		tok, err := dec.Token() // Token interface example of discriminated union
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement: // push
			stack = append(stack, tok.Name.Local)
		case xml.EndElement: // pop
			stack = stack[:len(stack)-1]
		case xml.CharData:
			if containsAll(stack, os.Args[1:]) {
				fmt.Printf("%s: %s\n", strings.Join(stack, " "), tok)
			}
		}
	}
}

// containsAll reports whether x contains the elements of y, in order.
func containsAll(x, y []string) bool {
	for len(y) <= len(x) {
		if len(y) == 0 {
			return true
		}
		if x[0] == y[0] {
			y = y[1:]
		}
		x = x[1:]
	}
	return false
}

/*

// Usage
$ go build $cwd/fetch
$ go build $cwd/xmlselect
$ fetch http://www.w3.org/TR/2006/REC-xml11-20060816 | ./xmlselect div div h2
html body div div h2: 1 Introduction
html body div div h2: 2 Documents
html body div div h2: 3 Logical Structures
html body div div h2: 4 Physical Structures
html body div div h2: 5 Conformance
...

*/
