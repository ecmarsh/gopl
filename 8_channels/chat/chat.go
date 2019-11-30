// chat lets several users broadcast textual messages to each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// job of main is to listen for and accept incoming network connections
// from clients. for each one, it creates a new handleConn goroutine.
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

// broadcaster's local variable `clients` records current set of
// connected clients. only information recorded is the identity of its
// outgoing message channel.
type client chan<- string // an outgoing message channel

// broadcaster listens on global entering and leaving channels
// for announcements of arriving and departing clients. when it receives
// one of these events, it updates the clients set and if the event
// was a departure, it closes the client's outgoing message channel.
// it also listens for events on global channel, where the client sends
// all its incoming messages and broadcasts this message upon receipt of global.
var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming messages
)

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

// handleConn creates a new outgoing message channel for its client
// and announces the arrival of this client to the broadcaster over `entering`.
// then it reads every line of text from the client, sending each line to the
// broadcaster over the global incoming message channel, and prefixes each
// message with the identity of its sender. When there is nothing more to read,
// it announces the departure of the client over the `leaving` chan and closes.
func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go ClientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + "has arrived"
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	// NOTE: ignoring portential errors from input.err()

	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

// ClientWriter is created by handleConn for each
// client that receives messages broadcast to the client's outgoing
// message channel and writes them to the client's network connection.
// the loop terminates when the broadcaster closes the channel
// after receiving a `leaving` notification.
func ClientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

/*
Note that the only variables that are shared by multiple gourtines are
channels and instances of net.Conn, both of which are concurrency safe.
See concurrency and shared variables chapter for more on implications
of sharing variables and concurrency safety.


USAGE

$ go build $GOPATH/src/path/to/chat
$ go build $GOPATH/src/path/to/netcat
$ ./chat 8
$ ./netcat
You are 127.0.0.1:PORT1
127.0.0.1:PORT2 has arrived
> Hi!
127.0.0.1:PORT1: Hi
127.0.0.1:PORT2 : Hi Yourself!

#  On another channel
./netcat
You are 127.0.0.1:PORT2
127.0.0.1:PORT1: Hi
> Hi Yourself!
127.0.0.1:PORT2 : Hi Yourself!

...
^C
*/
