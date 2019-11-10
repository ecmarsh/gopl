// Reverb echos different "volumes" of stdin and sends it to the server.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func reverb(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

// handleConn reads input and prints the servers response as a second goroutine.
func handleConn(c net.Conn) {
	input := bufio.NewScanner(c)
	for input.Scan() {
		// use goroutine because a real echo consists of
		// the composition of the three independent shouts
		go reverb(c, input.Text(), 1*time.Second)
	}
	// NOTE: ignoring potential errors from input.Err()
	c.Close()
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}

/*
// Usage
$ go build -o echo ./reverb
$ go build ./netcat.go
$ ./echo &
$ ./netcat
> Is anybody there?
	 IS ANYBODY THERE?
	 Is anybody there?
	 is anybody there?
> Yoo-hooo!
	 YOO-HOOO!
Is	 Yoo-hooo!
 anybod	 yoo-hooo!
y there?
	 IS ANYBODY THERE?
	 Is anybody there?
	 is anybody there?
> Testing 1,2,3
	 TESTING 1,2,3
> hello
	 Testing 1,2,3
	 HELLO
	 testing 1,2,3
	 hello
	 hello
^C

# cleanup
$ killall echo
$ rm echo netcat
*/
