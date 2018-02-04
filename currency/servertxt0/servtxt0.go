package main

import (
	"flag"
	"fmt"
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
// This version of the server uses TCP sockets (or UDS) to implement a simple
// text-based application-level protocol. There are no streaming strategy
// employed for the read/write operations. Buffers are read in one shot
// creating opportunities for missing data during read.
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

	if _, err := conn.Write([]byte("Connected...\nUsage: GET <currency, country, or code>\n")); err != nil {
		log.Println("error writing:", err)
		return
	}

	// loop to stay connected with client until client breaks connection
	for {
		// buffer for client command
		cmdLine := make([]byte, (1024 * 4))
		n, err := conn.Read(cmdLine)
		if n == 0 || err != nil {
			log.Println("connection read error:", err)
			return
		}
		cmd, param := parseCommand(string(cmdLine[0:n]))
		if cmd == "" {
			if _, err := conn.Write([]byte("Invalid command\n")); err != nil {
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
				if _, err := conn.Write([]byte("Nothing found\n")); err != nil {
					log.Println("failed to write:", err)
				}
				continue
			}
			// send each currency info as a line to the client wiht fmt.Fprintf()
			for _, cur := range result {
				_, err := conn.Write([]byte(
					fmt.Sprintf(
						"%s %s %s %s\n",
						cur.Name, cur.Code, cur.Number, cur.Country,
					),
				))
				if err != nil {
					log.Println("failed to write response:", err)
					return
				}
			}

		default:
			if _, err := conn.Write([]byte("Invalid command\n")); err != nil {
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
