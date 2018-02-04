package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	curr "github.com/vladimirvivien/go-networking/currency/lib"
)

var (
	currencies = curr.Load("../data.csv")
)

// This program implements a simple currency lookup service
// over TCP or Unix Data Socket. It loads ISO currency
// information using package curr (see above) and uses a simple
// JSON-encode text-based protocol to exchange data with a client.
//
// Clients send currency search requests as JSON objects
// as {"Get":"<currency name,code,or country"}. The request data is
// then unmarshalled to Go type curr.CurrencyRequest using
// the encoding/json package.
//
// The request is then used to search the list of
// currencies. The search result, a []curr.Currency, is marshalled
// as JSON array of objects and sent to the client.
//
// Focus:
// This version of the code continues to improve on the robustness of
// the server code by introducing configuration for read and write timeout
// values.  This ensures that a client cannot hold a connection hostage by
// taking a long time to send or receive data.
//
// Testing:
// Netcat can be used for rudimentary testing.  However, use clientjsonX
// programs functional tests.
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
		log.Println(err)
		os.Exit(1)
	}
	defer ln.Close()
	log.Println("**** Global Currency Service ***")
	log.Printf("Service started: (%s) %s\n", network, addr)

	// delay to sleep when accept fails with a temporary error
	acceptDelay := time.Millisecond * 10
	acceptCount := 0

	// connection loop
	for {
		conn, err := ln.Accept()
		if err != nil {
			switch e := err.(type) {
			case net.Error:
				// if temporary error, attempt to connect again
				if e.Temporary() {
					if acceptCount > 5 {
						log.Printf("unable to connect after %d retries: %v", err)
						return
					}
					acceptDelay *= 2
					acceptCount++
					time.Sleep(acceptDelay)
					continue
				}
			default:
				log.Println(err)
				conn.Close()
				continue
			}
			acceptDelay = time.Millisecond * 10
			acceptCount = 0
		}
		log.Println("Connected to ", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

// handle client connection
func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("error closing connection:", err)
		}
	}()

	// set initial deadline prior to entering
	// the client request/response loop to 45 seconds.
	// This means that the client has 45 seconds to send
	// its initial request or loose the connection.
	if err := conn.SetDeadline(time.Now().Add(time.Second * 45)); err != nil {
		log.Println("failed to set deadline:", err)
		return
	}

	// command-loop
	for {
		dec := json.NewDecoder(conn)
		var req curr.CurrencyRequest
		if err := dec.Decode(&req); err != nil {
			switch err := err.(type) {
			//network error: disconnect
			case net.Error:
				// is it a timeout error?
				// A deadline policy maybe implemented here using a decreasing
				// grace period that eventually causes an error if reached.
				// Here we just reject the connection if timeout is reached.
				if err.Timeout() {
					fmt.Println("deadline reached, disconnecting...")
				}
				fmt.Println("network error:", err)
				return
			default:
				if err == io.EOF {
					fmt.Println("closing connection:", err)
					return
				}
				enc := json.NewEncoder(conn)
				if encerr := enc.Encode(&curr.CurrencyError{Error: err.Error()}); encerr != nil {
					fmt.Println("failed error encoding:", encerr)
					return
				}
				continue
			}
		}

		// search currencies, result is []curr.Currency
		result := curr.Find(currencies, req.Get)

		// send result
		enc := json.NewEncoder(conn)
		if err := enc.Encode(&result); err != nil {
			switch err := err.(type) {
			case net.Error:
				fmt.Println("failed to send response:", err)
				return
			default:
				if encerr := enc.Encode(&curr.CurrencyError{Error: err.Error()}); encerr != nil {
					fmt.Println("failed to send error:", encerr)
					return
				}
				continue
			}
		}

		// renew deadline for 45 secs later
		if err := conn.SetDeadline(time.Now().Add(time.Second * 90)); err != nil {
			fmt.Println("failed to set deadline:", err)
			return
		}
	}
}
