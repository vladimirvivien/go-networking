package main

import (
	"flag"
	"fmt"
	"net"
)

var (
	host string
)

// this program looks up the IP addresses associated
// with the hostname.  It uses the default resolver
// which most likely will use a DNS lookup.
// To force the resolver to use Cgo, set the following
// environment variable:
// GODEBUG=netdns=cgo or
// GODEBUG=netdns=cgo+1
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
