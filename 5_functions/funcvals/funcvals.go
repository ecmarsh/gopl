// Package funcvals includes examples for function values.
package funcvals

import "fmt"

// Node represents an html Node for this example
type Node struct {
	Type                    NodeType
	Data                    string
	Attr                    []Attribute
	FirstChild, NextSibling *Node
}

// NodeType of html Node
type NodeType int32

// Specific html Node types
const (
	ErrorNode NodeType = iota
	TextNode
	DocumentNode
	ElementNode
	CommentNode
	DoctypeNode
)

// Attribute for html Node
type Attribute struct {
	Key, Val string
}

// forEachNode calls the functions pre(x) and post(x) for each node
// x in the tree rooted at n. Both functions are optional.
// pre is called before the children are visited (preorder) and
// post is called after (postorder).
func forEachNode(n *Node, pre, post func(n *Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

var depth int

func startElement(n *Node) {
	if n.Type == ElementNode {
		fmt.Printf("%*s<%s>\n", depth*2, "", n.Data)
		depth++
	}
}
func endElement(n *Node) {
	if n.Type == ElementNode {
		depth--
		fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
	}
}

// Note * in adverb %*s prints a string padded with
// variable number of spaces where the width and string
// are provided by arguments depth*2 and "".

// Usage: forEachNode(doc, startElement, endElement)
