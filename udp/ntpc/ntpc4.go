package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

// This program implements an NTP client that is capable of
// using either UDP or Unix Domain Socket datagram.  To do this,
// the program uses the Dialer to explicitly configure the client
// dialing process.
//
// The program uses -host to specify the remote address
// (or socket path) and -n for the network protocl ("udp" or "datagram").
func main() {
	var host string
	var network string
	flag.StringVar(&host, "e", "us.pool.ntp.org:123", "NTP host")
	flag.StringVar(&network, "n", "udp", "network protocol to use")
	flag.Parse()

	// req data packet is a 48-byte long value
	// that is used for sending time request.
	req := make([]byte, 48)

	// req is initialized with 0x1B or 0001 1011 which is
	// a request setting for time server.
	req[0] = 0x1B

	// rsp byte slice used to receive server response
	rsp := make([]byte, 48)

	// Create a Dialer which allows us to specify dialing options.
	// We will need this a bit later to configure the local address
	// when the program is using "unixgram"
	dialer := net.Dialer{}

	// IMPORTANT: when network is "unixgram", the local address
	// must be created and set explicitly (see ntpc2.go).
	if network == "unixgram" {
		laddr := &net.UnixAddr{Name: fmt.Sprintf("%s-client", host), Net: network}
		dialer.LocalAddr = laddr
	}

	// Setup connection (net.Conn) with Dial()
	conn, err := dialer.Dial(network, host)
	if err != nil {
		fmt.Printf("failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Println("failed while closing connection:", err)
		}
	}()

	// Once connection is established, the code pattern
	// is the same as in the previos impl.

	fmt.Printf("time from (%s) (%s)\n", network, conn.RemoteAddr())

	// send time request
	if _, err = conn.Write(req); err != nil {
		fmt.Printf("failed to send request: %v\n", err)
		os.Exit(1)
	}

	// block to receive server response
	read, err := conn.Read(rsp)
	if err != nil {
		fmt.Printf("failed to receive response: %v\n", err)
		os.Exit(1)
	}
	//ensure we read 48 bytes back (NTP protocol spec)
	if read != 48 {
		fmt.Println("did not get all expected bytes from server")
		os.Exit(1)
	}

	// NTP data comes in as big-endian (LSB [0...47] MSB)
	// with a 64-bit value containing the server time in seconds
	// where the first 32-bits are seconds and last 32-bit are fractional.
	// The following extracts the seconds from [0...[40:43]...47]
	// it is the number of secs since 1900 (NTP epoch)
	secs := binary.BigEndian.Uint32(rsp[40:])
	frac := binary.BigEndian.Uint32(rsp[44:])

	// Many OSs use Unix time epoch which is num of secs since 1970,
	// while NTP's ephoch starts on Jan 1, 1900.  Therefore,
	// to get the correct time, we must adjust the epocs properly
	// by removing 70 yrs of seconds (1970-1900) offset.
	ntpEpoch := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	unixEpoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	offset := unixEpoch.Sub(ntpEpoch).Seconds()
	now := float64(secs) - offset
	fmt.Printf("%v\n", time.Unix(int64(now), int64(frac)))
}
