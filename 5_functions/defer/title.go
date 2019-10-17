// Package main fetches an HTML document and prints title
package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"../funcvals"
	"golang.org/x/net/html"
)

// Attribute for html Node
type Attribute struct {
	Key, Val string
}

// Print HTML doc title or error if not Content-Type HTML
func titlenodefer(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	// Check Content-Type is HTML (text/html; charset=utf-8)
	ct := resp.Header.Get("Content-Type")
	if ct != "text/html" && !strings.HasPrefix(ct, "text/html;") {
		resp.Body.Close()
		return fmt.Errorf("%s has type %s, not text/html", url, ct)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" &&
			n.FirstChild != nil {
			fmt.Println(n.FirstChild.Data)
		}
	}
	funcvals.ForEachNode(doc, visitNode, nil)
	return nil
}

// Print HTML doc title or error if not Content-Type HTML
func titledefer(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // Executes after function completion

	ct := resp.Header.Get("Content-Type")
	if ct != "text/html" && !strings.HasPrefix(ct, "text/html;") {
		return fmt.Errorf("%s has type %s, not text/html", url, ct)
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" &&
			n.FirstChild != nil {
			fmt.Println(n.FirstChild.Data)
		}
	}
	funcvals.ForEachNode(doc, visitNode, nil)
	return nil
}

func main() {
	for _, arg := range os.Args[1:] {
		if err := titledefer(arg); err != nil {
			fmt.Fprintf(os.Stderr, "title: %v\n", err)
		}
	}
}
