# Vision

The Open Reserve Decentralized Web Node (DWN) is a Golang implementation of [DIF Decentralized Web Node Spec](https://identity.foundation/decentralized-web-node/spec/). The service is a part of a
larger [Decentralized Web](https://identity.foundation/) architecture. The DWN is a service that is used by its owner
to store and share data, and orchestrate data activities among other participants.

Goals include:

* To create a MVP of a DWN service, written in GoLang, that can be used in an "[as-a-service](https://forums.tbd.website/t/dwn-sdks-and-as-a-service/128)" capacity.  
* This DWN service should be able to be hosted by a domain owner (ie dwn.openreserve.io), and can support multiple wallets/devices/DIDs exchanging messages


# Guiding Principles

The service is assumed to be run by **a single organization** and assumes **external authentication and authorization**.
The service assumes no infrastructure requirements and is flexible to multiple deployment models, databases, key
management solutions, and user interfaces. We expect that a wide array of users and use cases will use and build on top
of the service, creating layers of abstraction and intermediation for processing business logic, handling user accounts,
and so on.

The service may choose to support both synchronous and asynchronous APIs; though generally, it should limit
statefulness. The service may implement a set of “static” APIs that expose functionality like storing data, 
sending and/or relaying messages, and other capabilities required for a decentralized web5 interaction

# Feature Support

The future feature set of the DWN is largely influenced by the standards and features specified in the [DIF Decentralized Web Node Spec](https://identity.foundation/decentralized-web-node/spec/),
in aim of advancing the adoption of decentralized web nodes "as a service." It adheres to best practices and guidelines for
implementing a privacy-minded, secure, and performant service.

We favor evaluating the addition of features and standards on a case-by-case basis, and looking towards implementations
of standards and features that are well-reasoned, with committed developers and use cases. Features that already
demonstrated usage and interoperability outside the project are prime candidates for adoption.

## Language Support

The DWN ecosystem uses a wide set of tools, languages, and technologies: working across web browsers, mobile
applications, backend servers, and more. This service
uses [Golang](https://go.dev/) because of its robust
cryptographic support, speed, ability to be compiled to [WASM](https://webassembly.org/), and, above all else,
simplicity. It is crucial that the code we write is approachable to encourage contribution. Simple and clear is always
preferred over clever.

The future is multi-language, and multi-platform. We welcome initiatives for improving multi-language and multi-platform
support, and are open to incubating them in our GitHub organization. When future SDKs are developed, it is expected that
they follow the same feature set and API as the Go SDK in addition to fulfilling the suite of language interoperability
tests.