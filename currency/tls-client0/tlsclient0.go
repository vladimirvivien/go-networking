package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	curr "github.com/vladimirvivien/go-networking/currency/lib"
)

const prompt = "currency"

// This porgram is a client implementation for the currency service
// program.  It sends JSON-encoded requests, i.e. {"Get":"USD"},
// and receives a JSON-encoded array of currency information directly
// over TCP or unix domain socket.
//
// Focus:
// This version of the client program highlights the use of
// IO streaming, data serialization, and client-side error handling.
//
// Usage: client [options]
// options:
//  - e service endpoint or socket path, default localhost:4443
//  - n network protocol name [tcp,unix], default tcp
//
// Once started a prompt is provided to interact with service.
func main() {
	// setup flags
	var addr, network, ca string
	flag.StringVar(&addr, "e", "localhost:4443", "service endpoint [ip addr or socket path]")
	flag.StringVar(&network, "n", "tcp", "network protocol [tcp,unix]")
	flag.StringVar(&ca, "ca", "../certs/ca-cert.pem", "CA certificate")
	flag.Parse()

	// Load our CA certificate
	caCert, err := ioutil.ReadFile(ca)
	if err != nil {
		log.Fatal("failed to read CA cert", err)
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	// TLS configuration
	tlsConf := &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            certPool,
	}

	// create a tls.Conn to connect to server
	conn, err := tls.Dial(network, addr, tlsConf)
	if err != nil {
		log.Fatal("failed to create socket:", err)
	}
	defer conn.Close()
	fmt.Println("connected to currency service: ", addr)

	var param string

	// start REPL
	for {
		fmt.Println("Enter search string or *")
		fmt.Print(prompt, "> ")
		_, err = fmt.Scanf("%s", &param)
		if err != nil {
			fmt.Println("Usage: <search string or *>")
			continue
		}

		req := curr.CurrencyRequest{Get: param}

		// Send request:
		// use json encoder to encode value of type curr.CurrencyRequest
		// and stream it to the server via net.Conn.
		if err := json.NewEncoder(conn).Encode(&req); err != nil {
			switch err := err.(type) {
			case net.Error:
				fmt.Println("failed to send request:", err)
				continue
			default:
				fmt.Println("failed to encode request:", err)
				continue
			}
		}

		// Receive response
		var currencies []curr.Currency
		err = json.NewDecoder(conn).Decode(&currencies)
		if err != nil {
			switch err := err.(type) {
			case net.Error:
				fmt.Println("failed to receive response:", err)
				continue
			default:
				fmt.Println("failed to decode response:", err)
				continue
			}
		}

		// print currencies
		for i, c := range currencies {
			fmt.Printf("%2d. %s[%s]\t%s, %s\n", i, c.Code, c.Number, c.Name, c.Country)
		}
	}
}
