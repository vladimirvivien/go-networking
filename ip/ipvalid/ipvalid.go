package main

import (
	"fmt"
	"net"
	"os"
)

// A program that validates IP addresses.
// The program accepts an IP address as command parameter.
// ./ipvalid <ip-address>
func main() {
	if len(os.Args) != 2 {
		return
	}

	// parse IP address
	ip := net.ParseIP(os.Args[1])

	// if not-nil IP is OK, else bad
	if ip != nil {
		fmt.Printf("%v OK\n", ip)
	} else {
		fmt.Println("bad address")
	}
}
