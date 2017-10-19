package main

import (
	"fmt"
	"net"
)

func main() {
	addr0 := "74.125.21.113:80"
	if ip, port, err := net.SplitHostPort(addr0); err == nil {
		fmt.Printf("ip=%s port=%s\n", ip, port)
	} else {
		fmt.Println(err)
	}

	addr1 := "[2607:f8b0:4002:c06::65]:80"
	if ip, port, err := net.SplitHostPort(addr1); err == nil {
		fmt.Printf("ip=%s port=%s\n", ip, port)
	} else {
		fmt.Println(err)
	}

	local := ":8080"
	if ip, port, err := net.SplitHostPort(local); err == nil {
		fmt.Printf("ip=%s port=%s\n", ip, port)
	} else {
		fmt.Println(err)
	}
}
