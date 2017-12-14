package main

import (
	"flag"
	"fmt"
	"math"
	"net"
	"os"
)

var (
	cidr string
)

func init() {
	flag.StringVar(&cidr, "c", "", "the CIDR address")
}

func main() {
	flag.Parse()

	if cidr == "" {
		fmt.Println("CIDR address missing")
		os.Exit(1)
	}

	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Println("failed parsing CIDR address: ", err)
		os.Exit(1)
	}

	ones, totalBits := ipnet.Mask.Size()
	size := totalBits - ones
	totalHosts := math.Pow(2, float64(size))
	wildcardIP := wildcard(net.IP(ipnet.Mask))
	last := lastIP(ip, net.IPMask(wildcardIP))

	fmt.Println()
	fmt.Printf("CIDR: %s\n", cidr)
	fmt.Println("------------------------------------------------")
	fmt.Printf("CIDR Block:     %s\n", cidr)
	fmt.Printf("Network:        %s\n", ipnet.IP)
	fmt.Printf("IP Range:       %s - %s\n", ip, last)
	fmt.Printf("Total Hosts:    %0.0f\n", totalHosts)
	fmt.Printf("Netmask:        %s\n", net.IP(ipnet.Mask))
	fmt.Printf("Wildcard Mask:  %s\n", wildcardIP)
	fmt.Println()
}

func wildcard(mask net.IP) net.IP {
	var ipVal net.IP
	for _, octet := range mask {
		ipVal = append(ipVal, ^octet)
	}
	return ipVal
}

func lastIP(ip net.IP, mask net.IPMask) net.IP {
	ipIn := ip.To4()
	if ipIn == nil {
		ipIn = ip.To16()
	}
	var ipVal net.IP
	for i, octet := range ipIn {
		ipVal = append(ipVal, octet|mask[i])
	}
	return ipVal
}
