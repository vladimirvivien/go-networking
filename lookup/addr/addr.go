package main

import (
	"context"
	"flag"
	"fmt"
	"net"
)

var (
	addr string
)

// this program performs a reverse lookup on the given
// host address.  This example uses the net.Resolver type
// directly to show how to programmatically configure the
// resolver to default to Go or Cgo.  To force it to use
// Cgo, set
//
//net.Resolver{PreferGo:false}
func main() {
	flag.StringVar(&addr, "addr", "127.0.0.1", "host address to lookup")
	flag.Parse()

	resolver := net.Resolver{PreferGo: true}
	names, err := resolver.LookupAddr(context.Background(), addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(names)
}
