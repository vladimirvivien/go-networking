package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

// This program implements a simple echo server over that is able
// to use TCP or Unix Domain Socket (streaming).
// When the server receives a request, it returns its content immediately.
//
// Usage:
// echos2
//   -e <endpoint: ip addr or path>
//   - n <protoco [tcp,unix]>
func main() {
	var addr string
	var network string
	flag.StringVar(&addr, "e", ":4040", "service endpoint [ip addr or socket path]")
	flag.StringVar(&network, "n", "tcp", "network protocol [tcp,unix]")
	flag.Parse()

	// validate network
	switch network {
	case "tcp", "tcp4", "tcp6", "unix":
	default:
		fmt.Println("unsupported network protocol")
		os.Exit(1)
	}

	// announce service using the Listen function
	// which creates a generic Listen listener.
	l, err := net.Listen(network, addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer l.Close()
	fmt.Printf("listening at (%s) %s\n", network, addr)

	// req/response loop
	for {
		// use Listener to block and wait for connection
		// request using function Accept() which then
		// creates a generic Conn value.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("failed to accept conn:", err)
			conn.Close()
			continue
		}
		fmt.Println("connected to: ", conn.RemoteAddr())

		go handleConnection(conn)
	}
}

// handleConnectino reads request from connection
// with conn.Read() then write response using
// conn.Write().  Then the connection is closed.
func handleConnection(conn net.Conn) {
	defer conn.Close() // clean up when done

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	// echo buffer
	w, err := conn.Write(buf[:n])
	if err != nil {
		fmt.Println("failed to write to client:", err)
		return
	}
	if w != n { // was all data sent
		fmt.Println("warning: not all data sent to client")
		return
	}
}
