// Crawl explores web links in parallel (concurrently). Note: termination not addressed.
package main

import (
	"fmt"
	"log"
	"os"
)

// main shows alternative solution to execessive concurrency.
// Uses sequential crawl, but calls from one of 20 long-lived crawler
// goroutines, ensuring that at most 20 HTTP requests are active.
// Crawler goroutines all fed by same channel, unseenLinks, and
// the main goroutine is responsible for de-duplicating items it
// receives from the worklist and sending each unseen over to new goroutine.
func main() {
	worklist := make(chan []string)  // list of URLs, may have duplicates
	unseenLinks := make(chan string) // de-duplicated URLs

	// Add CLI arguments to worklist.
	go func() { worklist <- os.Args[1:] }()

	// Create 20 crawler goroutines to fetch each unseen link.
	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link)
				go func() { worklist <- foundLinks }()
			}
		}()
	}

	// The main goroutine de-duplicates worklist items
	// and sends the unseen ones to the crawlers.
	// seen map is confined so it can only be accessed by current go routine
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				unseenLinks <- link
			}
		}
	}
}

// tokens is a counting semaphore, used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

// countingSemaphoreMain uses tokens to max routines to 20.
func countingSemaphoreMain() {
	worklist := make(chan []string)
	var n int // number of pending sends to worklist

	// Start with CLI args
	n++
	go func() { worklist <- os.Args[1:] }()

	// Crawl the web concurrently
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string) {
					worklist <- crawl(link)
				}(link)
			}
		}
	}
}

// Sending a goroutine into channel acquires a token, and
// releasing a gourtine from channel releases token, creating new vacancy.
func crawl(url string) []string {
	fmt.Println(url)
	tokens <- struct{}{} // acquire a token
	list, err := links.Extract(url)
	<-tokens // release the token
	if err != nil {
		log.Print(err)
	}
	return list
}

// main resembles BFS. A worklist records the queue of items that
// need processing, each item being a list of URLs to crawl, but
// instead of a queue, a channel is used. Each call to `crawl` occurs
// in its own goroutine and sends discovered links to the workload.
func sequentialMain() {
	worklist := make(chan []string)

	// Start with the CLI arguments.
	go func() { worklist <- os.Args[1:] }()

	// Crawl the web concurrently
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				go func(link string) {
					worklist <- sequentialCrawl(link)
				}(link)
			}
		}
	}
}

// Noncurrent crawl for refrence
func sequentialCrawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}
