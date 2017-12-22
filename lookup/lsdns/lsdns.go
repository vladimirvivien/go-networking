package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
)

var (
	ip    string
	host  string
	ns    bool
	mx    bool
	txt   bool
	cname bool
)

func init() {
	flag.StringVar(&ip, "ip", "", "IP address for DNS operation")
	flag.StringVar(&host, "host", "", "Host address for DNS operation")
	flag.BoolVar(&ns, "ns", false, "Host name server lookup")
	flag.BoolVar(&mx, "mx", false, "Host domain mail server lookup")
	flag.BoolVar(&txt, "txt", false, "Host domain TXT lookup")
	flag.BoolVar(&cname, "cname", false, "Host CNAME lookup")
}

type lsdns struct {
	resolver *net.Resolver
}

func newLsdns() *lsdns {
	return &lsdns{net.DefaultResolver}
}

func (ls *lsdns) reverseLkp(ip string) error {
	names, err := ls.resolver.LookupAddr(context.Background(), ip)
	if err != nil {
		return err
	}

	fmt.Println("Reverse lookup")
	fmt.Println("--------------")
	for _, name := range names {
		fmt.Println(name)
	}
	return nil
}

func (ls *lsdns) hostLkp(host string) error {
	addrs, err := ls.resolver.LookupHost(context.Background(), host)
	if err != nil {
		return err
	}

	fmt.Println("Host lookup")
	fmt.Println("-----------")
	for _, addr := range addrs {
		fmt.Printf("%-30s%-20s\n", host, addr)
	}
	return nil
}

func (ls *lsdns) nsLkp(host string) error {
	nses, err := ls.resolver.LookupNS(context.Background(), host)
	if err != nil {
		return err
	}
	fmt.Println("NS lookup")
	fmt.Println("---------")
	for _, ns := range nses {
		fmt.Printf("%-25s%-20s\n", host, ns.Host)
	}
	return nil
}

func (ls *lsdns) mxLkp(host string) error {
	mxes, err := ls.resolver.LookupMX(context.Background(), host)
	if err != nil {
		return err
	}
	fmt.Println("MX lookup")
	fmt.Println("---------")
	for _, mx := range mxes {
		fmt.Printf("%-17s%-11s\n", host, mx.Host)
	}
	return nil
}

func (ls *lsdns) txtLkp(host string) error {
	txts, err := ls.resolver.LookupTXT(context.Background(), host)
	if err != nil {
		return err
	}
	fmt.Println("TXT lookup")
	fmt.Println("---------")
	for _, txt := range txts {
		fmt.Printf("%-17s%-11s\n", host, txt)
	}
	return nil
}

func (ls *lsdns) cnameLkp(host string) error {
	name, err := ls.resolver.LookupCNAME(context.Background(), host)
	if err != nil {
		return err
	}
	fmt.Println("CNAME lookup")
	fmt.Println("------------")
	fmt.Printf("%s: %s\n", host, name)
	return nil
}

func main() {
	flag.Parse()
	ls := newLsdns()

	switch {
	case ip != "":
		ls.reverseLkp(ip)
	case host != "":
		switch {
		case ns:
			ls.nsLkp(host)
		case mx:
			ls.mxLkp(host)
		case txt:
			ls.txtLkp(host)
		case cname:
			ls.cnameLkp(host)
		default:
			ls.hostLkp(host)
		}
	default:
		fmt.Println("flag ip or host must be provided")
		os.Exit(1)
	}
}
