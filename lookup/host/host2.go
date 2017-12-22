package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
)

var (
	host string
)

// this program looks up the IP addresses associated
// with the hostname.  This program uses the net.Resolver
// type directly to specify how the resolver works
// net.Resolver{PureGo:true}
func main() {
	flag.StringVar(&host, "host", "localhost", "host name to resolve")
	flag.Parse()

	res := net.Resolver{PreferGo: true}
	addrs, err := res.LookupHost(context.Background(), host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(addrs)
}
