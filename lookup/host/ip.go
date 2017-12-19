package main

import (
	"context"
	"flag"
	"fmt"
	"net"
)

var (
	host string
)

// this program looks up the IP addresses associated
// with the hostname using LookupIP which is similar to LookupHost.
// It uses the default resolver
// which most likely will use a DNS lookup.
// To force the resolver to use Cgo, set the following
// environment variable:
// GODEBUG=netdns=cgo
func main() {
	flag.StringVar(&host, "host", "localhost", "host name to resolve")
	flag.Parse()

	res := net.Resolver{PreferGo: true}
	addrs, err := res.LookupIPAddr(context.Background(), host)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(addrs)
}
