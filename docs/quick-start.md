# Quick Start Guide

## Start a Local Node for Development and Testing

### Prerequisites

* Docker
* Docker Compose

### Docker Compose Instructions

```shell
git clone https://github.com/openreserveio/dwn.git
cd dwn/deployment/localdev/compose
docker compose -f ./full-deploy.yaml build
docker compose -f ./full-deploy.yaml up
```

You should be able to hit port :8080 on your localhost to access the DWN service.
