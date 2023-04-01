# Quick Start Guide

## Start a Local Node for Development and Testing using Docker Compose

### Prerequisites

* Docker
* Docker Compose

### Docker Compose Instructions

```sh
git clone https://github.com/openreserveio/dwn.git
cd dwn/deployment/localdev/compose
docker compose -f ./full-deploy.yaml build
docker compose -f ./full-deploy.yaml up
```

You should be able to hit port :8080 on your localhost to access the DWN service.

## Build DWN on your local machine and run the individual components

### Prerequisites

* Golang v1.20.2 or higher
* Make

### Build Instructions

```sh
# Clone Repo
git clone https://github.com/openreserveio/dwn.git
cd dwn/go

# Build Executable
make build-executable
```

The output of the build will be in `go/build/release`, and you can use that to run any of the [individual services](req_arch_design/arch/services.md) of the DWN node.

