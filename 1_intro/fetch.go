// Prints content found at each specified URL.
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, url := range os.Args[1:] {
		if !strings.HasPrefix(url, "http") {
			pre := []string{"http://", url}
			url = strings.Join(pre, "")
		}
		res, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		if _, err := io.Copy(os.Stdout, res.Body); err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %s: %v\n", url, err)
			os.Exit(1)
		}
		res.Body.Close()
	}
}
