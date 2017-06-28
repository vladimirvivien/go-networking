# Go Network Programming
This repository contains notes and samples source code for the video series, of the same name, published by Packt Publishing. 

## A quick network review
Before we jump headfirst into writing networked programs in Go, let us do a quick review of the concepts we will cover in this video series.  To be clear, when we say *network*, in the contex of this discussion, we are referring to connected computers on a network communicating using protocols such UDP or TCP over IP.  Given the scope of this session, we are not going to have deep discussion about these protocols as it is assume that you have some familiarity with them.  However, it is worth reviewing the three main protocols that will impact our discussions for the remainder of this session.

 * IP - the ability of computers to communicate with other computers on the same network is specified by the Internet Protocol (or IP).  This protocol specifies concepts such as Addressing and Routing to enable it to send data packets (datagrams) from one host to another host connected on the same network or across network boundaries.  In the video session, we will explore the Go types and functions available in the net package to support IP.
 * UDP -  If IP specifies how computers are networked, the User Datagram Protocol (or UDP) was one of the earlier protocols built on top of IP that uses a connectionless transmission approach to send data from one host to another.  We will see how to use the types provided by the Go API to work with UDP.
 * TCP - When the transmission of data, between hosts on a network, requires a more robust guarantees of delivery than, say UDP, the Transmission Control Protocol (or TCP) is used over an IP network.  In this session we will see how the Go 


## The Net Package
When writing programs that communicate over a network in Go, you will likely start with the *net* package (https://golang.org/pkg/net).  This package and its sub-packages provide a rich API that expose low-level networking primitives as well as application-level protocols. 



Each logical component of a network is represented by a type including hardware interfaces, network, packets, address, protocols, and connections.  Furthermore, each type exposes a multitude of methods giving Go one of the most complete standard libraries for network programming supporting both IPv4 and IPv6.
