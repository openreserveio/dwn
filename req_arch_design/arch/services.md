# Service Design

## Goals

* Lightweight services architecture
* Easy to docker compose up and get going (both for devs and other users)
* API services serve API consumers and should be lightweight orchestrators to purpose-specific Backend services
* Backend services are modular gRPC components which serve a specific purpose, adhere to most DDD design principles, and stay lightweight where possible
* Easy docker compose up!

## API Services


### HTTP/API Service

* Exposes a client-facing API that adheres to the [DIF Decentralized Web Node Spec](https://identity.foundation/decentralized-web-node/spec)
* Basic validations but should be focused on message validation, and orchestration logic. Use message authentication service to authenticate messages vs. all connections
* Scale out, generally stateless

## Backend Services

* Document DB (mongodb) for storage
* Cache service (?TBD) for caching of signer info, possibly
* CollectionsService for management of collections, collection schema definitions -> docdb
* KeyService for message authentication, message signing, key management -> docdb
