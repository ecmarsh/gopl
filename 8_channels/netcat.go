// Netcat is a read-only TCP client.
package main

import (
	"io"
	"log"
	"net"
	"os"
)

// The program reads data from the connection and writes it to stdout
// until an EOF condition or error occurs.
func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mustCopy(os.Stdout, conn)
}

// mustCopy is a utility used in many examples to handle error writing.
func mustCopy(dest io.Writer, src io.Reader) {
	if _, err := io.Copy(dest, src); err != nil {
		log.Fatal(err)
	}
}
