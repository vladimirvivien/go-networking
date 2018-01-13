package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

// This program implements a simple echo server over TCP.
// When the server receives a request, it returns its content
// immediately.
//
// Usage:
// echos -e <host:address>
func main() {
	var addr string
	flag.StringVar(&addr, "e", ":4040", "service address endpoint")
	flag.Parse()

	// create local addr for socket
	laddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// announce service using ListenTCP
	// which a TCPListener.
	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("listening at (tcp)", laddr.String())

	// req/response loop
	for {
		// use TCPListener to block and wait for TCP
		// connection request using AcceptTCP which creates a TCPConn
		conn, err := l.AcceptTCP()
		if err != nil {
			fmt.Println("failed to accept conn:", err)
			conn.Close()
			continue
		}
		fmt.Println("connected to: ", conn.RemoteAddr())

		go handleConnection(conn)
	}
}

// handleConnection reads request from connection
// with conn.Read() then write response using
// conn.Write().  Then the connection is closed.
func handleConnection(conn *net.TCPConn) {
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
