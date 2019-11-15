// du reports the disk usage of one or more directories specified via CLI.
// Most of its work is done by the `walkdir` function which enumerates entries
// using the entries of the directory `dir` using the `dirents` helper function.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
func walkDir(dir string, fileSizes chan<- int64) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkDir(subdir, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

// dirents returns the entries of directory `dir`.
func dirents(dir string) []os.FileInfo {
	// ReadDir returns a slice of os.FileInfo, the same info
	// that a call to os.Stat returns for a single file.
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}

// main1 uses two gourtines: the background goroutine calls `walkDir` for
// each directory specified on the CLI and finally closes the fileSizes channel.
// The main goroutine computes the sum of the file sizes it receives from the channel and prints total.
// Note it does not print its progress, and if we simply moved printDU into loop, would print a lot of output.
func main1() {
	// Determine the initial directories
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	// Traverse the file tree.
	fileSizes := make(chan int64)
	go func() {
		for _, root := range roots {
			walkDir(root, fileSizes)
		}
		close(fileSizes)
	}()

	// Print the results.
	var nfiles, nbytes int64
	for size := range fileSizes {
		nfiles++
		nbytes += size
	}
	printDiskUsage(nfiles, nbytes)
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

// main2 keeps us informed of progress. It prints the totals periodically but only
// if the -v  flag is specified since not all users will want to see progress msgs.
// The bg goroutine that loops over `roots` is unchanged, but the main goroutine
// uses a ticker to generate events every 500ms and a selct statement to wait for
// either the a file size message (where it updates totals) or a tick event,
// where it prints the current totals. With no -v flag, the select statement is effectively disabled.
var verbose = flag.Bool("v", false, "show verbose progress messages")
func main2() {
	// ...start background goroutine...
	// Determine the initial directories
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	fileSizes := make(chan int64)
	go func() {
		for _, root := range roots {
			walkDir(root, fileSizes)
		}
		close(fileSizes)
	}()

	// Print the results periodically.
	var tick <-chan time.Time
	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}
	var nfiles, nbytes int64
loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if !ok {
				break loop // fileSizes was closed, break for and `loop`.
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes) // final totals
}

// The final improvement speeds up traversal. There's no reason why all calls
// to walkDir can't be done currently, exploiting parallelism in disk system,
// so we can create a new gourtine for each call to `walkDir` by utilizing
// `sync.WaitGroup` to count the number of calls to walkDir,
// and closes `fileSizes` channel when counter drops to zero.
// Lastly, since final main may create thousands of gouroutines at peak,
// we must change dirents to use a counting semaphore to prevent it from
// opening too many files at once, just as done in concurrent web crawler, ./crawl.
func main() {
	// ... determine roots (initial file directories) ...
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Traverse each root of the file tree in parallel.
	fileSizes := make(chan int64)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go parallelWalkDir(root, &n, fileSizes)
	}
	go func() {
		n.Wait()
		close(fileSizes)
	}()

	// ...select loop...
	// Print the results periodically.
	var tick <-chan time.Time
	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}
	var nfiles, nbytes int64
loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if !ok {
				break loop // fileSizes was closed, break for and `loop`.
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes) // final totals
}

func parallelWalkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go parallelWalkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

// sema is a counting semaphore to for limiting concurrency in dirents.
var sema = make(chan struct{}, 20)

// semaDirents returns the entries of directory dir.
func semaDirents(dir string) []os.FileInfo {
	sema <- struct{}{}  	   // acquire token
	defer func() { <- sema }() // release token

	// ... same as unlimited dirents above ...
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}

/*
// USAGE
$ go build $GOPATH/path/to/du
$ ./du -v $HOME /usr /bin /etc
##### files    #.# GB
##### files    #.# GB
..........     ......
 */
