SEMVER := 0.0.7-alpha
PLATFORMS := linux/amd64,linux/arm64

generate-protobuf:
	protoc --proto_path=../servicedef --go_out=generated/services/ --go_opt=paths=source_relative --go-grpc_out=generated/services/ --go-grpc_opt=paths=source_relative ../servicedef/backend_services.proto
	protoc --proto_path=../servicedef --go_out=generated/events/ --go_opt=paths=source_relative --go-grpc_out=generated/events/ --go-grpc_opt=paths=source_relative ../servicedef/events.proto

generate-mocks: generate-protobuf
	mockgen -source storage/collection_store.go -destination generated/mocks/collection_store.go -package mocks
	mockgen -source generated/services/backend_services_grpc.pb.go -destination generated/mocks/backend_services_grpc.pb.go -package mocks

buildx-docker-use:
	docker buildx create --use

build-docker-base: buildx-docker-use
	docker buildx build --platform=$(PLATFORMS) -f applications/pkg/docker/Dockerfile-base -t openreserveio/base .
	docker tag openreserveio/base:latest ghcr.io/openreserveio/base:latest
	docker tag openreserveio/base:latest ghcr.io/openreserveio/base:$(SEMVER)
	docker push ghcr.io/openreserveio/base:latest
	docker push ghcr.io/openreserveio/base:$(SEMVER)

build-docker-api: build-docker-base
	docker buildx build --platform=$(PLATFORMS) -f applications/pkg/docker/Dockerfile-api -t openreserveio/api .

build-docker-collsvc: build-docker-base
	docker buildx build --platform=$(PLATFORMS) -f applications/pkg/docker/Dockerfile-collsvc -t openreserveio/collsvc .

build-docker-keysvc: build-docker-base
	docker buildx build --platform=$(PLATFORMS) -f applications/pkg/docker/Dockerfile-keysvc -t openreserveio/keysvc .

build-docker-all: build-docker-base build-docker-api build-docker-collsvc build-docker-keysvc

tag-docker-github-packages: build-docker-all
	docker tag openreserveio/api:latest ghcr.io/openreserveio/api:latest
	docker tag openreserveio/api:latest ghcr.io/openreserveio/api:$(SEMVER)
	docker tag openreserveio/collsvc:latest ghcr.io/openreserveio/collsvc:latest
	docker tag openreserveio/collsvc:latest ghcr.io/openreserveio/collsvc:$(SEMVER)
	docker tag openreserveio/keysvc:latest ghcr.io/openreserveio/keysvc:latest
	docker tag openreserveio/keysvc:latest ghcr.io/openreserveio/keysvc:$(SEMVER)


push-docker-github-packages: tag-docker-github-packages
	docker push ghcr.io/openreserveio/api:$(SEMVER)
	docker push ghcr.io/openreserveio/api:latest
	docker push ghcr.io/openreserveio/collsvc:$(SEMVER)
	docker push ghcr.io/openreserveio/collsvc:latest
	docker push ghcr.io/openreserveio/keysvc:$(SEMVER)
	docker push ghcr.io/openreserveio/keysvc:latest

build-release-executables:
	mkdir -p build/release
	cd applications/dwn && GOOS=darwin GOARCH=arm64 go build -o ../../build/release/dwn-darwin-arm64-$(SEMVER) .
	cd applications/dwn && GOOS=darwin GOARCH=amd64 go build -o ../../build/release/dwn-darwin-amd64-$(SEMVER) .
	cd applications/dwn && GOOS=linux GOARCH=amd64 go build -o ../../build/release/dwn-linux-amd64-$(SEMVER) .
	cd applications/dwn && GOOS=linux GOARCH=arm64 go build -o ../../build/release/dwn-linux-arm64-$(SEMVER) .
	cd applications/dwn && GOOS=linux GOARCH=386 go build -o ../../build/release/dwn-linux-386-$(SEMVER) .
	cd applications/dwn && GOOS=windows GOARCH=amd64 go build -o ../../build/release/dwn-windows-amd64-$(SEMVER) .
	cd applications/dwn && GOOS=windows GOARCH=386 go build -o ../../build/release/dwn-windows-386-$(SEMVER) .