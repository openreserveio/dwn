# DWN

GoLang-based implementation of [DIF's DWN Specification](https://identity.foundation/decentralized-web-node/spec/), following [TBD](https://developer.tbd.website)'s efforts in the space.


## Goals

* To create a MVP of a DWN service, written in GoLang, that can be used in an "[as-a-service](https://forums.tbd.website/t/dwn-sdks-and-as-a-service/128)" capacity.  
* This DWN service should be able to be hosted by a domain owner (ie dwn.openreserve.io), and can support multiple wallets/devices/DIDs exchanging messages

## Progress

### For initial v0.1.0 release (2022 Q4)

- [x] [Initial API layer](https://github.com/openreserveio/dwn/pull/4)
- [x] [Initial message processing, collections write/read](https://github.com/openreserveio/dwn/pull/9)
- [x] [Create Record and Initial Entry](https://github.com/openreserveio/dwn/pull/19)
- [x] [CollectionsWrite](https://github.com/openreserveio/dwn/issues/23), [CollectionsDelete](https://github.com/openreserveio/dwn/issues/26), [CollectionsCommit](https://github.com/openreserveio/dwn/issues/25) Message Chains
- [ ] [CollectionsQuery with filters](https://github.com/openreserveio/dwn/issues/27)
- [ ] [Hooks](https://github.com/openreserveio/dwn/issues/28) and [Webhooks](https://github.com/openreserveio/dwn/issues/29), [Events](https://github.com/openreserveio/dwn/issues/12)
- [ ] [Permissions](https://github.com/openreserveio/dwn/issues/30)
- [x] Docker-based deployment focus

### For v0.2.0 release (2023 Q1)

- [ ] Protocols support
- [ ] Sync
- [ ] Commit Strategy changes to message processing
- [ ] DID service for openreserve DID method
- [ ] Support [DIDComm](https://identity.foundation/didcomm-messaging/spec/)

### To be scheduled:

- [ ] Support [DID Registration](https://identity.foundation/did-registration/)


# Documentation

* [Quick Start Guide](docs/quick-start.md)
* [Developer Documentation](docs/developer/README.md)
* [Operations Documentation](docs/operations/README.md)

