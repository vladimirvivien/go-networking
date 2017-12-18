package main

import (
	"fmt"
	"net"
	"os"
)

// Easiest way to create net.IP value is to use
// net.ParseIP which parses a string value representation
// of a IPv4 dot-separated or IPv6 colon-separated address.
// This example uses net.ParseIP to parse an IP address provided
// from the command line and prints information about the address.
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Missing ip address")
		os.Exit(1)
	}
	ip := net.ParseIP(os.Args[1])
	if ip == nil {
		fmt.Println("Unable to parse IP address.")
		fmt.Println("Address should use IPv4 dot-notation or IPv6 colon-notation")
		os.Exit(1)
	}

	fmt.Println()
	fmt.Printf("IP:             %s\n", ip)
	fmt.Printf("Default Mask:   %s\n", net.IP(ip.DefaultMask()))
	fmt.Printf("Loopback:       %t\n", ip.IsLoopback())
	fmt.Println("Unicast:")
	fmt.Printf("  Global:       %t\n", ip.IsGlobalUnicast())
	fmt.Printf("  Link:         %t\n", ip.IsLinkLocalUnicast())
	fmt.Println("Multicast:")
	fmt.Printf("  Global:       %t\n", ip.IsMulticast())
	fmt.Printf("  Interface     %t\n", ip.IsInterfaceLocalMulticast())
	fmt.Printf("  Link          %t\n", ip.IsLinkLocalMulticast())
	fmt.Println()
}
