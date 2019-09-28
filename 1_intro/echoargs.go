// Prints its CLI arguments
package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	var s, sep string
	// TODO: Benchmark (time pkg) appending string vs `strings.Join`
	for i, arg := range os.Args[1:] {
		n := strconv.FormatInt(int64(i+1), 10)
		s += sep + "(" + n + ")" + arg
		sep = " "
	}
	fmt.Println(s)
}
