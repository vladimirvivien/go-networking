package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

// This program implements a simple echo client over tcp.
// It sends a text content to the server and displays
// the response on the screen.
//
// Usage: echoc -e <host-addr-endpoint> <text content>
func main() {
	var addr string
	flag.StringVar(&addr, "e", "localhost:4040", "service address endpoint")
	flag.Parse()
	text := flag.Arg(0)

	// use ResolveTCPAddr to create address to connect to
	raddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Use DialTCP to create a connection to the
	// remote address. Note that there is no need
	// to specify the local address.
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		fmt.Println("failed to connect to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// send text to server
	_, err = conn.Write([]byte(text))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// read response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("failed reading response:", err)
		os.Exit(1)
	}
	fmt.Println(string(buf[:n]))
}
