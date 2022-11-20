SEMVER := 0.0.2-alpha

generate-protobuf:
	protoc --proto_path=../servicedef --go_out=generated/services/ --go_opt=paths=source_relative --go-grpc_out=generated/services/ --go-grpc_opt=paths=source_relative ../servicedef/backend_services.proto


build-docker-base:
	docker build -f applications/pkg/docker/Dockerfile-base -t openreserveio/base .

build-docker-api: build-docker-base
	docker build -f applications/pkg/docker/Dockerfile-api -t openreserveio/api .

build-docker-collsvc: build-docker-base
	docker build -f applications/pkg/docker/Dockerfile-collsvc -t openreserveio/collsvc .

build-docker-all: build-docker-base build-docker-api build-docker-collsvc

tag-docker-github-packages: build-docker-all
	docker tag openreserveio/api:latest ghcr.io/openreserveio/api:latest
	docker tag openreserveio/api:latest ghcr.io/openreserveio/api:$(SEMVER)
	docker tag openreserveio/collsvc:latest ghcr.io/openreserveio/collsvc:latest
	docker tag openreserveio/collsvc:latest ghcr.io/openreserveio/collsvc:$(SEMVER)


push-docker-github-packages: tag-docker-github-packages
	docker push ghcr.io/openreserveio/api:$(SEMVER)
	docker push ghcr.io/openreserveio/api:latest
	docker push ghcr.io/openreserveio/collsvc:$(SEMVER)
	docker push ghcr.io/openreserveio/collsvc:latest


