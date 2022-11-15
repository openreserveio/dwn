SEMVER := 0.0.1-alpha

build-docker-base:
	docker build -f applications/pkg/docker/Dockerfile-base -t openreserveio/base .

build-docker-api: build-docker-base
	docker build -f applications/pkg/docker/Dockerfile-api -t openreserveio/api .

build-docker-all: build-docker-base build-docker-api

tag-docker-github-packages: build-docker-all
	docker tag openreserveio/api:latest ghcr.io/openreserveio/api:latest
	docker tag openreserveio/api:latest ghcr.io/openreserveio/api:$(SEMVER)

push-docker-github-packages: tag-docker-github-packages
	docker push ghcr.io/openreserveio/api:$(SEMVER)
	docker push ghcr.io/openreserveio/api:latest


