package main

import (
	"flag"
	"fmt"
	"net"
)

var (
	host string
)

// this program looksup the ip address
// for a given host address
func main() {
	flag.StringVar(&host, "host", "localhost", "host name to resolve")
	flag.Parse()

	addrs, err := net.LookupHost(host)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(addrs)
}
