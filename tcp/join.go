package main

import (
	"fmt"
	"net"
)

func main() {
	host := "10.10.100.1"
	port := "1234"

	addr := net.JoinHostPort(host, port)

	fmt.Println("addr = ", addr)
}
