# Go Network Programming
This repository contains notes and samples source code for the video series, of the same name, published by Packt Publishing. 

## A quick network review
Before we jump headfirst into writing networked programs in Go, let us do a quick review of the concepts we will cover in this video series.  To be clear, when we say *network*, in the context of this discussion, we are referring to connected computers on a network communicating using protocols such UDP or TCP over IP.  Given the scope of this session, we are not going to have deep discussion about these protocols as it is assume that you have some familiarity with them.  However, it is worth reviewing the three main protocols that will impact our discussions for the remainder of this session.

 * IP - the ability of computers to communicate with other computers on the same network is specified by the Internet Protocol (or IP).  This connectionless protocol specifies concepts such as Addressing and Routing to enable it to send data packets (datagrams) from one host to another host connected on the same network or across network boundaries.  In the video session, we will explore the Go types and functions available in the net package to support IP.
 
 * UDP - The User Datagram Protocol (UDP) is a core tenant of the Internet protocols.  It is, what is known as, a connectionless transmission protocol designed to reduce latency of delivery.  To do this, UDP sends data packets (or datagrams) to hosts, on an IP network, with minimum reliability and delivery guarantee.  In this session, we will explore how to use the constructs provided by the Go API to work with UDP.
 
 * TCP - When the transmission of data, between hosts on a network, requires a more robust guarantees of delivery than, say UDP, the Transmission Control Protocol (or TCP) is used over an IP network.  The protocol uses a session between communicating parties which starts with a handshake to establish a durable connection between hosts that can handle transmission error, out-of-order packets, and delivery guarantee.  In this session we will explore how Go supports TCP and the types available in the API to work with the protocol.


## The net Package
As mentioned in the opening, when writing programs that communicate over a network in Go, you will likely start with the *net* package (https://golang.org/pkg/net).  This package, and its sub-packages, provide a rich API that exposes low-level networking primitives as well as application-level protocols such as HTTP.  For this discussion, we will focus on protocols such IP, UDP, and TCP.

All logical components that makes up network communications can be  represented using types from the net package including:

- *Addressing* [address and host resolution in `net`]
- *Protocols* (IP, IP/TCP, IP/UDP) [protocol representation in `net`]
- *Network IO* [interfaces available for network communications]

Furthermore, each API interface exposes a multitude of methods giving Go one of the most complete standard libraries for Internet programming supporting both `IPv4` and `IPv6`. 

So now, let us take a look at the major API themes found in the net package.

### Addressing
One of the most basic primitives, when doing network programming, is the address.  It is used to identify networks and nodes interconnected together.  Many types and functions in the  `net` package support the representations of `IP` addresses using string literals as shown in the following:

```sh
"127.0.0.1"
"2607:f8b0:4002:c06::65"
```
As shown in the previous snippet, the `net` package supports the representation of both `IPv4` and `IPv6`.  Addresses can also include a service port when associated with protocols such as TCP/IP or UDP/IP.  For instance, the following string literal represents a host address along with a HTTP service accessible on port 80:

```sh
"74.125.21.113:80"
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
Similarly, functions and methods that work with `UDP` and `TCP` protocols directly have their own rich representation of addresses respectively.  For instance, type `UDPAddr` is shown below.
```go
type UDPAddr struct {
        IP   IP
        Port int
        Zone string 
}
```
Likewise, type `TCPAddr` is offered as a way for representing TCP/IP addresses as shown below.
```go
type IPAddr struct {
        IP   IP
        Port int
        Zone string 
}
```

### Protocols
Many functions and other types, in the net packages, use string literals to identify the protocol used when communicating across the network.  The following lists the string values of network protocols that we will be discussing in this series:
```sh
"ip",  "ip4",  "ip6"
"tcp", "tcp4", "tcp6" 
"udp", "udp4", "udp6"                        
```
The suffix `4` indicates a protocol for IPv4 only and the suffix `6` indicates a protocol using IPv6 only.  When the string literal omits the version, it targets IPv4 by default.  

For instance, the following calls function `Dial` (discussed later) providing it with the network protocol name as `"tcp"` and an address.

```go
Dial("tcp", "74.125.21.113:80")
```
When the network is `ip`, the string must include a protocol number (or name), separated by a colon.  For instance, the following function call specifies an IP network with the ICMP protocol.
```go
ResolveIPAddr("ip:icpm", "192.168.1.136")
```


# Topics 

- *Network interface information* [accessing hardware interface info]
- Address Resolution [lookup and resolving addresses]
- IP Communication
- UDP Communication
- TCP Communication
