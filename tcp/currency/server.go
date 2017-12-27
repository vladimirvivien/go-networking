package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"

	curr "github.com/vladimirvivien/go-networking/tcp/curlib"
)

var (
	currencies = curr.Load("./data.csv")
)

// This program implements a simple currency lookup service
// over TCP or Unix Data Socket. It loads ISO currency
// information using package curlib (see above) and makes
// and serves it using JSON-enocoded data.
//
// Clients send currency search requests as JSON objects such
// as {"Get":"USD"}. The request data is then unmarshalled to Go
// type curr.CurrencyRequest{Get:"USD"} using the encoding/json
// package.
//
// The request is then used to search the list of
// currencies. The search result, a []curr.Currency, is marshalled
// to JSON array of objects and send to the client.
//
// Data Serialization:
// This version of the program highlights the use of
// encoding to serialize data to/from Go data types
// to another representation (JSON).
//
// Usage: server [options]
// options:
//   -e host endpoint, default ":4040"
//   -n network protocol [tcp,unix], default "tcp"
func main() {
	// setup flags
	var addr string
	var network string
	flag.StringVar(&addr, "e", ":4040", "service endpoint [ip addr or socket path]")
	flag.StringVar(&network, "n", "tcp", "network protocol [tcp,unix]")
	flag.Parse()

	// validate supported network protocols
	switch network {
	case "tcp", "tcp4", "tcp6", "unix":
	default:
		fmt.Println("unsupported network protocol")
		os.Exit(1)
	}

	// create a listener for provided network and host address
	ln, err := net.Listen(network, addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer ln.Close()
	fmt.Println("**** Global Currency Service ***")
	fmt.Printf("Service started: (%s) %s\n", network, addr)

	// connection loop
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		fmt.Println("Connected to ", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

// handle client connection
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// loop to keep connection alive until client breaks connection
	for {
		// block to read in request in 4k buffer
		buf := make([]byte, 1024*4)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("error reading:", err)
			return
		}

		// unmarshal request into value of type curr.CurrencyRequest
		// Note unmarshal uses a sub-slice buf[:n]
		var req curr.CurrencyRequest
		if err := json.Unmarshal(buf[:n], &req); err != nil {
			fmt.Println("failed to unmarshal request:", err)
			// inform the client of bad request by marshaling
			// curr.CurrencyError to JSON to the client.
			cerr, jerr := json.Marshal(&curr.CurrencyError{Error: err.Error()})
			if jerr != nil {
				fmt.Println("failed to marshal CurrencyError:", jerr)
				continue
			}

			if _, werr := conn.Write(cerr); werr != nil {
				fmt.Println("failed to write to CurrencyError:", werr)
				return
			}
			continue
		}

		// search currencies, result is []curr.Currency
		result := curr.Find(currencies, req.Get)

		// marshal result to JSON array
		rsp, err := json.Marshal(&result)
		if err != nil {
			fmt.Println("failed to marshal data:", err)
			continue
		}

		// write response to client
		if _, err := conn.Write(rsp); err != nil {
			fmt.Println("failed to write response:", err)
			return
		}
	}
}
