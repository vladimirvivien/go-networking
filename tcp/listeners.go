package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("udp", "127.0.0.1:4040")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%T\n", l)
}
