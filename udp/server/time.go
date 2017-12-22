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
	var host string
	flag.StringVar(&host, "host", ":1123", "server address")
	flag.Parse()

	// create a UDP address
	addr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// setup UDP socket
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("failed to create socket:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("listening for time request: %s (%s)\n", addr, conn.LocalAddr())

	// read incoming request, but throw it away
	_, target, err := conn.ReadFrom(make([]byte, 48))
	if err != nil {
		fmt.Println("error getting request:", err)
		os.Exit(1)
	}

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
	if _, err := conn.WriteTo(rsp, target); err != nil {
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
