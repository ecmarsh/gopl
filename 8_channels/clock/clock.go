// Clock is a TCP server that periodally writes the time.
package main

import (
	"io"
	"log"
	"net"
	"time"
)

// The concurrent clock server version is implemented by
// adding the `go` keyword to the call to handleConn,
// which causes each call to run as its own goroutine.
// It performs the same function as seqClock, except now
// multiple clients can receive the time at once.
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn) // handle connection concurrently
	}
}

// seqClock is a sequential clock server.
func seqClock() {
	// Listen creates a net.Listener, which is an object for
	// incoming connections on a network port (in this case TCP 8000).
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
		handleConn(conn) // handle one connection at a time
	}
}

// handleConn handles one complete client connection.
// In a loop, it writes the current time, time.Now(),
// to the client using the io.WriteString interface, satisfied by net.Conn.
// The loop ends when write fails, likely due to client disconnect.
// After failure, deferred close is called to close the connection.
func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

/*
// Usage
$ cd 8_channels
$ go build -o clockserver ./clock
$ ./clockserver &
$ go build ./netcat.go
$ ./netcat
HH:MM:SS
HH:MM:SS
...
# Open new terminal
$ ./netcat
HH:MM:SS
HH:MM:SS (same as above)
...
$ killall netcat
$ kill %1 # kill concclock

# alternatively run both in the background
$ ./netcat >> log1 &
$ ./netcat >> log2 &
$ killall netcat
$ diff log1 log2 # output shows some lines same
$ rm clockserver netcat log* # cleanup
*/
