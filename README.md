# Go Network Programming
This repository contains notes and samples source code for the video series, of the same name, published by Packt Publishing. 

## A quick network review
Before we jump headfirst into writing networked programs in Go, let us do a quick review of the concepts we will cover in this video series.  To be clear, when we say *network*, in the context of this discussion, we are referring to connected computers on a network communicating using protocols such UDP or TCP over IP.  Given the scope of this session, we are not going to have deep discussion about these protocols as it is assume that you have some familiarity with them.  However, it is worth reviewing the three main protocols that will impact our discussions for the remainder of this session.

 * IP - the ability of computers to communicate with other computers on the same network is specified by the Internet Protocol (or IP).  This connectionless protocol specifies concepts such as Addressing and Routing to enable it to send data packets (datagrams) from one host to another host connected on the same network or across network boundaries.  In the video session, we will explore the Go types and functions available in the net package to support IP.
 
 * UDP - The User Datagram Protocol (UDP) is a core tenant of the Internet protocols.  It is, what is known as, a connectionless transmission protocol designed to reduce latency of delivery.  To do this, UDP sends data packets (or datagrams) to hosts, on an IP network, with minimum reliability and delivery guarantee.  In this session, we will explore how to use the constructs provided by the Go API to work with UDP.
 
 * TCP - When the transmission of data, between hosts on a network, requires a more robust guarantees of delivery than, say UDP, the Transmission Control Protocol (or TCP) is used over an IP network.  The protocol uses a session between communicating parties which starts with a handshake to establish a durable connection between hosts that can handle transmission error, out-of-order packets, and delivery guarantee.  In this session we will explore how Go supports TCP and the types available in the API to work with the protocol.


## The net Package
As mentioned in the opening, when writing programs that communicate over a network in Go, you will likely start with the *net* package (https://golang.org/pkg/net).  This package and its sub-packages provide a rich API that exposes low-level networking primitives as well as application-level protocols such as HTTP.  For this discussion, we will focus on protocols such IP, UDP, and TCP.

### Understanding the net package
All logical components that take part in network communications can be  represented by in Go including `hardware interfaces`, `network`, `packets`, `address`, `protocols`, and `host connections`.  Furthermore, each type exposes a multitude of methods giving Go one of the most complete standard libraries for Internet programming supporting both `IPv4` and `IPv6`.
