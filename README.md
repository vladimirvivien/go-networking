# Go Network Programming
This repository contains notes and samples source code for the video series, of the same name, published by Packt Publishing. 

## A quick network review
Before we jump headfirst into writing networked programs in Go, let us do a quick review of the concepts we will cover in this video series.  To be clear, when we say *network*, in the contex of this discussion, we are referring to a connected computers on a network communicating using protocols TCP or UDP over IP.  We are not going to have deep discussion about the background of these protocols as it is assume that you have some familiarity with them.  However, it is worth revisiting the followingn concepts as it will make some of our discussions that we cover in later topics clearer.

 * Internetworking 
 * Address
 * Session Protocol: TCP
 * Sessionless Protocol


## The Net Package
When writing programs that communicate over a network in Go, you will likely start with the *net* package (https://golang.org/pkg/net).  This package and its sub-packages provide a rich API that expose low-level networking primitives as well as application-level protocols. 



Each logical component of a network is represented by a type including hardware interfaces, network, packets, address, protocols, and connections.  Furthermore, each type exposes a multitude of methods giving Go one of the most complete standard libraries for network programming supporting both IPv4 and IPv6.
