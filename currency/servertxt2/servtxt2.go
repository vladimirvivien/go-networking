package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	curr "github.com/vladimirvivien/go-networking/currency/lib0"
)

var currencies = curr.Load("../data.csv")

// This program implements a simple currency lookup service
// over TCP or Unix Data Socket. It loads ISO currency
// information using package lib (see above) and uses a simple
// text-based protocol to interact with the client and send
// the data.
//
// Clients send currency search requests as a textual command in the form:
//
// GET <currency, country, or code>
//
// When the server receives the request, it is parsed and is then used
// to search the list of currencies. The search result is then printed
// line-by-line back to the client.
//
// Focus:
// This version of the currency server focuses on implementing a streaming
// strategy when receiving data from client to avoid dropping data when the
// request is larger than the internal buffer. This version uses the bufio
// package to use buffered readers to stream from net.Conn.
//
// Testing:
// Netcat or telnet can be used to test this server by connecting and
// sending command using the format described above.
//
// Usage: server0 [options]
// options:
//   -e host endpoint, default ":4040"
//   -n network protocol [tcp,unix], default "tcp"
func main() {
	var addr string
	var network string
	flag.StringVar(&addr, "e", ":4040", "service endpoint [ip addr or socket path]")
	flag.StringVar(&network, "n", "tcp", "network protocol [tcp,unix]")
	flag.Parse()

	// validate supported network protocols
	switch network {
	case "tcp", "tcp4", "tcp6", "unix":
	default:
		log.Fatalln("unsupported network protocol:", network)
	}

	// create a listener for provided network and host address
	ln, err := net.Listen(network, addr)
	if err != nil {
		log.Fatal("failed to create listener:", err)
	}
	defer ln.Close()
	log.Println("**** Global Currency Service ***")
	log.Printf("Service started: (%s) %s\n", network, addr)

	// connection-loop - handle incoming requests
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			if err := conn.Close(); err != nil {
				log.Println("failed to close listener:", err)
			}
			continue
		}
		log.Println("Connected to", conn.RemoteAddr())

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("error closing connection:", err)
		}
	}()

	if _, err := fmt.Fprint(conn, "Connected...\nUsage: GET <currency, country, or code>\n"); err != nil {
		log.Println("error writing:", err)
		return
	}

	// buffered reader to stream data using 4-byte chunks until ('\n\')
	// The chunks are kept small to demonstrate streaming using io.Reader.
	reader := bufio.NewReaderSize(conn, 4)

	// command-loop
	for {
		cmdLine, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println("connection read error:", err)
				return
			}
		}
		reader.Reset(conn) //cleans buffer

		cmd, param := parseCommand(cmdLine)
		if cmd == "" {
			if _, err := fmt.Fprint(conn, "Invalid command\n"); err != nil {
				log.Println("failed to write:", err)
				return
			}
			continue
		}

		// execute command
		switch strings.ToUpper(cmd) {
		case "GET":
			result := curr.Find(currencies, param)
			if len(result) == 0 {
				if _, err := fmt.Fprint(conn, "Nothing found\n"); err != nil {
					log.Println("failed to write:", err)
				}
				continue
			}
			// send each currency info as a line to the client wiht fmt.Fprintf()
			for _, cur := range result {
				_, err := fmt.Fprintf(
					conn,
					"%s %s %s %s\n",
					cur.Name, cur.Code, cur.Number, cur.Country,
				)

				if err != nil {
					log.Println("failed to write response:", err)
					return
				}
			}

		default:
			if _, err := fmt.Fprintf(conn, "Invalid command\n"); err != nil {
				log.Println("failed to write:", err)
				return
			}
		}
	}
}

func parseCommand(cmdLine string) (cmd, param string) {
	parts := strings.Split(cmdLine, " ")
	if len(parts) != 2 {
		return "", ""
	}
	cmd = strings.TrimSpace(parts[0])
	param = strings.TrimSpace(parts[1])
	return
}
