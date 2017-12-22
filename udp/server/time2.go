package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

// This program is a super simple network time protocol server.
// It uses UDP to return the number of seconds since 1900.
func main() {
	var port string
	flag.StringVar(&port, "p", "1123", "port value")
	flag.Parse()

	// Creaets a UDP address
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%s", port))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// get a socket, announce service on network
	// Because it's connectionless, the same socket can be
	// reused to handle with multiple request/responses.
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("failed to create socket:", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Printf("listening on (udp) %s\n", conn.LocalAddr())

	for {
		// block to read incoming requests
		_, target, err := conn.ReadFrom(make([]byte, 48))
		if err != nil {
			fmt.Println("error getting request:", err)
			os.Exit(1)
		}

		// handle request
		go handleRequest(conn, target)
	}
}

// handle incoming requests
func handleRequest(conn *net.UDPConn, addr net.Addr) {
	// get seconds and fractional secs since 1900
	secs, fracs := getNTPSeconds(time.Now())

	// response packet is filled with the seconds and
	// fractional sec values using Big-Endian
	rsp := make([]byte, 48)
	// write seconds (as uint32) in buffer at [40:43]
	binary.BigEndian.PutUint32(rsp[40:], uint32(secs))
	// write seconds (as uint32) in buffer at [44:47]
	binary.BigEndian.PutUint32(rsp[44:], uint32(fracs))

	// send data
	if _, err := conn.WriteTo(rsp, addr); err != nil {
		fmt.Println("err sending data:", err)
		os.Exit(1)
	}
}

// getNTPSecs decompose current time as NTP seconds
func getNTPSeconds(t time.Time) (int64, int64) {
	// convert time to total # of secs since 1970
	// add NTP epoch offets as total #secs between 1900-1970
	secs := t.Unix() + int64(getNTPOffset())
	fracs := t.Nanosecond()
	return secs, int64(fracs)
}

// getNTPOffset returns the 70yrs between Unix epoch
// and NTP epoch (1970-1900) in seconds
func getNTPOffset() float64 {
	ntpEpoch := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	unixEpoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	offset := unixEpoch.Sub(ntpEpoch).Seconds()
	return offset
}
