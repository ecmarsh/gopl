// Package main exemplifies a simple fulfilment of the http.Handler interface.
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	db := database{"shoes": 50, "socks": 5}
	log.Fatal(http.ListenAndServe("localhost:8000", db))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

func (db database) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

/*
// Usage
$ go build PATHFROMWORKSPACEROOT/http1
$ ./http1 & # start up the sever in background
$ ./fetch http://localhost:8000 # use fetch program to output db items
// Output
shoes: $50.00
socks: $5.00
*/
