# DWN

[![godoc dwn](https://img.shields.io/badge/godoc-dwn-blue)](https://pkg.go.dev/github.com/openreserveio/dwn/go)
[![go version 1.20.2](https://img.shields.io/badge/go_version-1.20.2-brightgreen)](https://go.dev/)
[![license Apache 2](https://img.shields.io/badge/license-Apache%202-black)](https://github.com/openreserveio/dwn/blob/main/LICENSE)
[![issues](https://img.shields.io/github/issues/openreserveio/dwn)](https://github.com/openreserveio/dwn/issues)
![push](https://github.com/openreserveio/dwn/workflows/dwn-compile-test-golang/badge.svg?branch=main&event=push)

GoLang-based implementation of [DIF's DWN Specification](https://identity.foundation/decentralized-web-node/spec/), following [TBD](https://developer.tbd.website)'s efforts in the space.


# Documentation

* [Quick Start Guide](docs/quick-start.md)
* [Developer Documentation](docs/developer/README.md)
* [Operations Documentation](docs/operations/README.md)

# Project Goals

* To create a MVP of a DWN service, written in GoLang, that can be used in an "[as-a-service](https://forums.tbd.website/t/dwn-sdks-and-as-a-service/128)" capacity.  
* This DWN service should be able to be hosted by a domain owner (ie dwn.openreserve.io), and can support multiple wallets/devices/DIDs exchanging messages

## Capabilities

### The v0.1.x releases includes the following functionality:

- [Initial API layer](https://github.com/openreserveio/dwn/pull/4)
- [Initial message processing, collections write/read](https://github.com/openreserveio/dwn/pull/9)
- [Create Record and Initial Entry](https://github.com/openreserveio/dwn/pull/19)
- [RecordWrite](https://github.com/openreserveio/dwn/issues/23), [RecordDelete](https://github.com/openreserveio/dwn/issues/26), [RecordCommit](https://github.com/openreserveio/dwn/issues/25) Message Chains
- [RecordQuery with filters](https://github.com/openreserveio/dwn/issues/27)
- [Public Alpha Deployment](https://github.com/openreserveio/dwn/issues/41)
- [Opinionated Client Library](https://github.com/openreserveio/dwn/issues/39)
- [Hooks](https://github.com/openreserveio/dwn/issues/28) and [Webhooks](https://github.com/openreserveio/dwn/issues/29), [Events](https://github.com/openreserveio/dwn/issues/12)
- Docker-based deployment focus

### Planned For v0.2.0 release (2023 Q3)

- [X] [Permissions](https://github.com/openreserveio/dwn/issues/30)
- [ ] [HooksDelete](https://github.com/openreserveio/dwn/issues/78)
- [ ] [ProtocolsConfigure](https://github.com/openreserveio/dwn/issues/79)
- [ ] [ProtocolsQuery](https://github.com/openreserveio/dwn/issues/80)
- [ ] Sync
- [X] Commit Strategy changes to message processing
- [ ] DID service for openreserve DID method
- [ ] Support [DIDComm](https://identity.foundation/didcomm-messaging/spec/)

### To be scheduled:

- [ ] Enforce [Protocols](https://identity.foundation/decentralized-web-node/spec/#protocols)
- [ ] Support [DID Registration](https://identity.foundation/did-registration/)
- [ ] Automated release to alpha.openreserve.io open testing channel




# Project Resources

| Resource                              | Description                                                                   |
|--------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------|
| [VISION](docs/VISION.md)              | Outlines the project vision                                                   |
| [DIPs](docs/dip/README.md)            | Proposals for improving the DWN components                                    |
| [VERSIONING](docs/VERSIONING.md)      | Project versioning strategy                                                   |
| [CODEOWNERS](CODEOWNERS)              | Outlines the project lead(s)                                                  |
| [CODE_OF_CONDUCT](CODE_OF_CONDUCT.md) | Expected behavior for project contributors, promoting a welcoming environment |
| [CONTRIBUTING](CONTRIBUTING.md)       | Developer guide to build, test, run, access CI, chat, discuss, file issues    |
| [GOVERNANCE](GOVERNANCE.md)           | Project governance                                                            |
| [SECURITY](SECURITY.md)               | Vulnerability and bug reporting                                               |
| [LICENSE](LICENSE)                    | Apache License, Version 2.0                                                   |
