// Package main drives the tempconv with Celsius Flag.
package main

import (
	"flag"
	"fmt"

	"../tempconv"
)

var temp = tempconv.CelsiusFlag("temp", 20.0, "the temperature")

func main() {
	flag.Parse()
	fmt.Println(*temp)
}

/*

// Usage

$ go build ./tempflag
$ ./tmpflag
20°C
$ ./tempflag -temp -18C
-18°C
$ ./tempflag -temp 212°F
100°C
$ ./tempflag -help
Usage of ./tempflag:
	-temp value
		the temperature (default 20°C)

*/
