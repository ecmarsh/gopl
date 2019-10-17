// Package soletitle exemplified how you might check panic and recover
package soletitle

import (
	"fmt"

	"../funcvals"
	"golang.org/x/net/html"
)

// soleTitle returns the text of the first non-empty title element
// in doc, and an error if there was not exactly one.
func soleTitle(doc *html.Node) (title string, err error) {
	type bailout struct{}

	defer func() {
		switch p := recover(); p {
		case nil:
			// no panic
		case bailout{}:
			// "expected" panic for multile "<title>" elements
			err = fmt.Errorf("multiple title elements")
		default:
			panic(p) // unexpected panic; let it panic
		}
	}()

	// Bail out of recursion if we find more than one non-empty title.
	funcvals.ForEachNode(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" &&
			n.FirstChild != nil {
			if title != "" {
				panic(bailout{}) // multiple title elements
			}
			title = n.FirstChild.Data
		}
	}, nil)

	// Normal functionality
	if title == "" {
		return "", fmt.Errorf("no title element")
	}
	return title, nil
}
