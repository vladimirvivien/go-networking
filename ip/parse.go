package main

import (
	"fmt"
	"net"
)

func main() {
	localIP := net.ParseIP("127.0.0.1")
	remoteIP := net.ParseIP("2607:f8b0:4002:c06::65")

	fmt.Println("local IP: ", localIP)
	fmt.Println("remote IP: ", remoteIP)
}
