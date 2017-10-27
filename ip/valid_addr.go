package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		return
	}
	addr := os.Args[1]
	ip := net.ParseIP(addr)
	if ip != nil {
		fmt.Printf("%v OK\n", ip)
	} else {
		fmt.Println("bad address")
	}
}
