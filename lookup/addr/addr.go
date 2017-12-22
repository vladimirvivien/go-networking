package main

import (
	"flag"
	"fmt"
	"net"
)

var (
	addr string
)

// this program performs a reverse lookup on the given
// host address.
//
func main() {
	flag.StringVar(&addr, "addr", "127.0.0.1", "host address to lookup")
	flag.Parse()

	names, err := net.LookupAddr(addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(names)
}
