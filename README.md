# Go Network Programming
This repository contains notes and samples source code for the video series, of the same name, published by Packt Publishing. 

## A quick network review
Before we jump headfirst into writing networked programs in Go, let us do a quick review of the concepts we will cover in this video series.  To be clear, when we say *network*, in the context of this discussion, we are referring to connected computers on a network communicating using protocols such UDP or TCP over IP.  Given the scope of this session, we are not going to have deep discussion about these protocols as it is assume that you have some familiarity with them.  However, it is worth reviewing the three main protocols that will impact our discussions for the remainder of this session.

 * *IP* - the ability of computers to communicate with other computers on the same network is specified by the Internet Protocol (or IP).  This connectionless protocol specifies concepts such as Addressing and Routing to enable it to send data packets (datagrams) from one host to another host connected on the same network or across network boundaries.  In the video session, we will explore the Go types and functions available in the net package to support IP.
 
 * *UDP* - The User Datagram Protocol (UDP) is a core tenant of the Internet protocols.  It is, what is known as, a connectionless transmission protocol designed to reduce latency of delivery.  To do this, UDP sends data packets (or datagrams) to hosts, on an IP network, with minimum reliability and delivery guarantee.  In this session, we will explore how to use the constructs provided by the Go API to work with UDP.
 
 * *TCP* - When the transmission of data, between hosts on a network, requires a more robust guarantees of delivery than, say UDP, the Transmission Control Protocol (or TCP) is used over an IP network.  The protocol uses a session between communicating parties which starts with a handshake to establish a durable connection between hosts that can handle transmission error, out-of-order packets, and delivery guarantee.  In this session we will explore how Go supports TCP and the types available in the API to work with the protocol.


## The net Package
As mentioned in the opening, when writing programs that communicate over a network in Go, you will likely start with the *net* package (https://golang.org/pkg/net).  This package, and its sub-packages, provide a rich API that exposes low-level networking primitives as well as application-level protocols such as HTTP.  For this discussion, we will focus on protocols such IP, UDP, and TCP.

Before we dive head-first into our discussion, it is worth taking a high-level look at the `net` package.  There are some critical themes represented in the package that should be discussed before we dive into the details.  For instance, all logical components that makes up network communications are abstracted as types and supporting functions.  Let us take a look at some of these.

### Addressing
One of the most basic primitives, when doing IP-based network programming, is the address.  Addresses are used to identify networks and network nodes interconnected together.  In the `net` package `IP` addresses can be represented using string literals with the dot notation for IPv4 and colon-delimited for IPv6 addresses as shown below:
```go
var localIP = "127.0.0.1"
var remIP = "2607:f8b0:4002:c06::65"
```
When dealing with an UDP or TCP, the address can also include a port number separated by a colon.  For IPv6 addresses, the IP address is enclosed within a bracket then followed by the port.
```go
var webAddr = "127.0.0.1:8080"
var sshAddr = "[2607:f8b0:4002:c06::65]:22"
```
As we explore the protocols in detail, we will see how each have their own typed representation of addresses such as `net.IPAddr`, `net.UDPAddr`, and `net.TCPAddr`.

### Name and service resolution
One crucial function of a network API is the ability to resolve services, addresses, and names from a given network.  The net package provides several functions to query naming and service information such as host names, IP, reverse lookup, NS, MX, and CNAME records.

For instance the following program looks up the IP address for the given host. It uses
function `net.LookupHost()` which returns a slice of string IP addresses.  

```go
func main() {
	flag.StringVar(&host, "host", "localhost", "host name to resolve")
	flag.Parse()

	addrs, err := net.LookupHost(host)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(addrs)
}
```
The net package uses `resolver strategy` to determined how resolve network names 
depending on the operatng system. By default, the resolver will attempt to use a pure Go 
mechanism that queries DNS directly to avoid OS-related thread penalties.  Given certain
conditions and combinations, the resolver may fallback to using a C-implemented resolver
to relies on OS system calls.

This behavior can be overridden using the `GODEBUG` environment variale as shown.

```sh
export GODEBUG=netdns=go    # use Go resolver
export GODEBUG=netdns=cgo   # use C resolver
```



### Protocols
Another prominent theme in the net package is `protocol` representation.  Many functions and types, in the `net` package, use string literals to identify the protocol that is being used when communicating across the network.  The following lists the string identifier for the protocols that we will cover:
```sh
"ip",  "ip4",  "ip6"
"tcp", "tcp4", "tcp6" 
"udp", "udp4", "udp6"                        
```
The suffix `4` indicates a protocol for IPv4 only and the suffix `6` indicates a protocol using IPv6 only.  When the string literal omits the version, it targets IPv4 by default.  We see this used many times during the video series when invoking functions and methods.

### Network communication
When building networked programs, you will certainly need a way to connect nodes together so they can communicate to exchange data. Depending on the nature of the network protocol you may also need: 

- the ability to announce the service on an available port
- the ability to listen, accept, and handle incoming client connections

Let us look at how the `net` package provides support for creating programs that can communicate on the network:

* *net.Conn* - this interface represents communication between two nodes on a network.  When writing networked application, eventually you will use an implementation of that interface to exchange data.  The net package comes with several implementations including `net.IPConn`, `net.UDPConn`, and `net.TCPConn` for low-level and protocol-specific functionalities.  For instance, streaming protocols such as TCP, exposes streaming IO semantics from `io.Reader` and `io.Writer` interfaces.  We will see how this is done as we get deeper into our sessions.

* *Listening for connections* - the `net` package provides several functions that can be used to listen for incoming connections depending on the protocol that you want to use.  These functions include `net.ListenIP()` and `net.ListenUDP()` which return a `net.IPConn` and `net.UDPConn` respectively.  To listen for TCP connections, we use `net.ListenTCP()` which returns `net.Listener` implementation.  

* *Dialing a connection* - to establish a connection from one network node to another, the `net` package uses the notion of dialing a connection provided by functions `net.IPDial()`, `net.UDPDial()`, and `net.TCPDial()` which return their respective connection implementations of `net.IPConn`, `net.UDPConn`, and `net.TCPConn`.  The wrapper function:
```go
func Dial(network, address string) (Conn, error)
```
is often used to create a connection.  It automatically returns the proper `net.Conn` implementation based on the network protocol specified in parameter `network`.  For instance, the following will open a TCP connection to the indicated address and port:
```go
net.Dial("tcp", "64.233.177.102")
```
We will see more examples of the Dial functions as we continue the video series.

* *Unix socket* - Lastly for completeness, it is worth mentioning that the `net` package also supports *Unix domain socket* protocol for doing both streaming and packet based inter-process communications.  The protocols are identified as "unix", "unixgram", and "unixpacket" and uses type `net.UnixConn` to represent a connection.  This video series does not discuss detail about these protocols.  However, they use similar interfaces and follow the same idioms as the TCP and UDP protocol implementations.  This should make them fairly easy to learn and use. 

## IP
Here is a n example of how IP addresses can be represented as string literals.  In the string literal value, the IPv4 uses dot notation to separate address bytes while IPv6 uses colon separator.  In both call to `net.ParseIP()`, the function parses the address and return a typed representation of the address or `nil` if the address is invalid.

```go
func main() {
	localIP := net.ParseIP("127.0.0.1")
	remoteIP := net.ParseIP("2607:f8b0:4002:c06::65")

	fmt.Println("local IP: ", localIP)
	fmt.Println("remote IP: ", remoteIP)
}
```
When the addresses is for a service associated with TCP or UDP (which will be covered in later sessions), the string literal can also include a service port.  For IPv4 addresses, the port is separated by colon and IPv6 addresses are placed within brackets followed by a colon and the port.  For instance, the following string literal represents a host address along with a HTTP service accessible on port 80.  Also when the address is assumed to be the local host, the IP address can be omitted, leaving only the colon followed by the port.  All three representations are valid versions of string representations of IP addresses in Go.

```go
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
```
Formally, the `net` package uses type `net.IP` to represent an IP address as a slice of bytes capable of storing both IPv4 and IPv6 addresses.  

```go
type IP []byte
```
The IP type exposes several methods that makes it easy to work and manipulate IP addresses. To illustrate this, the following code is a utility that validates then converts a given address to both IPv4 and IPv6 values.

```go
<sample code>
```
The working with the IP protocol directly, the net package in Go uses type `IPAddr` to provide a richer representation of an IP address which is used in several functions and methods.

```go
type IPAddr struct {
        IP   IP
        Zone string 
}
```


# Topics 

- *Network interface information* [accessing hardware interface info]
- Address Resolution [lookup and resolving addresses]
- IP Communication
- UDP Communication
- TCP Communication
