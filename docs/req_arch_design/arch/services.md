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
* Queueing system (NATS, ActiveMQ) for events
* [HookService](hook_service.md) for DWN Hook storage, management, and emitting of events to backend system to deliver hooks
  * HTTPS WebHook
  * gRPC Callback
  * WebSocket Message Send
  * APNS/FCM iPhone/Android
* Cache service (?TBD) for caching of signer info, possibly
* RecordService for management of records, data schema definitions -> docdb
* KeyService for message authentication, message signing, key management -> docdb
* NotificationService for delivering hook events
* API Service to export the standard [DWN HTTP interfaces](https://identity.foundation/decentralized-web-node/spec/#interfaces) to clients, as well as support Open Reserve's extensions.
* Relay Service to assist clients with DWN message relaying (to keep sender IP info confidential, for example)

