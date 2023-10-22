SEMVER := 0.2.0-alpha
PLATFORMS := linux/amd64,linux/arm64

generate-oapi:
	oapi-codegen -old-config-style -generate gin,types,client,spec -package notifycallback ../servicedef/notifycallback-oas.yaml >generated/oas/notifycallback/notifycallback.go


generate-protobuf:
	protoc --proto_path=../servicedef --go_out=generated/services/ --go_opt=paths=source_relative --go-grpc_out=generated/services/ --go-grpc_opt=paths=source_relative ../servicedef/backend_services.proto
	protoc --proto_path=../servicedef --go_out=generated/events/ --go_opt=paths=source_relative --go-grpc_out=generated/events/ --go-grpc_opt=paths=source_relative ../servicedef/events.proto

generate-mocks: generate-protobuf
	mockgen -source storage/record_store.go -destination generated/mocks/record_store.go -package mocks
	mockgen -source storage/hook_store.go -destination generated/mocks/hook_store.go -package mocks
	mockgen -source generated/services/backend_services_grpc.pb.go -destination generated/mocks/backend_services_grpc.pb.go -package mocks
	mockgen -source generated/events/events.pb.go -destination generated/mocks/events.pb.go -package mocks

build-docker-api:
	docker build -f applications/pkg/docker/Dockerfile-api -t openreserveio/api .

build-docker-recordsvc:
	docker build -f applications/pkg/docker/Dockerfile-recordsvc -t openreserveio/recordsvc .

build-docker-keysvc:
	docker build -f applications/pkg/docker/Dockerfile-keysvc -t openreserveio/keysvc .

build-docker-hooksvc:
	docker build -f applications/pkg/docker/Dockerfile-hooksvc -t openreserveio/hooksvc .

build-docker-notificationsvc:
	docker build -f applications/pkg/docker/Dockerfile-notificationsvc -t openreserveio/notificationsvc .

build-docker-all: build-docker-api build-docker-recordsvc build-docker-keysvc build-docker-hooksvc build-docker-notificationsvc

tag-docker-github-packages: build-docker-all
	docker tag openreserveio/api:latest ghcr.io/openreserveio/api:latest
	docker tag openreserveio/api:latest ghcr.io/openreserveio/api:$(SEMVER)
	docker tag openreserveio/recordsvc:latest ghcr.io/openreserveio/recordsvc:latest
	docker tag openreserveio/recordsvc:latest ghcr.io/openreserveio/recordsvc:$(SEMVER)
	docker tag openreserveio/keysvc:latest ghcr.io/openreserveio/keysvc:latest
	docker tag openreserveio/keysvc:latest ghcr.io/openreserveio/keysvc:$(SEMVER)
	docker tag openreserveio/hooksvc:latest ghcr.io/openreserveio/hooksvc:latest
	docker tag openreserveio/hooksvc:latest ghcr.io/openreserveio/hooksvc:$(SEMVER)
	docker tag openreserveio/notificationsvc:latest ghcr.io/openreserveio/notificationsvc:latest
	docker tag openreserveio/notificationsvc:latest ghcr.io/openreserveio/notificationsvc:$(SEMVER)


push-docker-multiplatform:

	docker buildx create --use
	docker buildx build --platform=$(PLATFORMS) --push -f applications/pkg/docker/Dockerfile-api -t ghcr.io/openreserveio/api:latest -t ghcr.io/openreserveio/api:$(SEMVER) .
	docker buildx build --platform=$(PLATFORMS) --push -f applications/pkg/docker/Dockerfile-recordsvc -t ghcr.io/openreserveio/recordsvc:latest -t ghcr.io/openreserveio/recordsvc:$(SEMVER) .
	docker buildx build --platform=$(PLATFORMS) --push -f applications/pkg/docker/Dockerfile-keysvc -t ghcr.io/openreserveio/keysvc:latest -t ghcr.io/openreserveio/keysvc:$(SEMVER) .
	docker buildx build --platform=$(PLATFORMS) --push -f applications/pkg/docker/Dockerfile-hooksvc -t ghcr.io/openreserveio/hooksvc:latest -t ghcr.io/openreserveio/hooksvc:$(SEMVER) .
	docker buildx build --platform=$(PLATFORMS) --push -f applications/pkg/docker/Dockerfile-notificationsvc -t ghcr.io/openreserveio/notificationsvc:latest -t ghcr.io/openreserveio/notificationsvc:$(SEMVER) .


build-executable:
	mkdir -p build/release
	cd applications/dwn && go build -o ../../build/release/dwn-$(SEMVER) .

build-release-executables:
	mkdir -p build/release
	cd applications/dwn && GOOS=darwin GOARCH=arm64 go build -o ../../build/release/dwn-darwin-arm64-$(SEMVER) .
	cd applications/dwn && GOOS=darwin GOARCH=amd64 go build -o ../../build/release/dwn-darwin-amd64-$(SEMVER) .
	cd applications/dwn && GOOS=linux GOARCH=amd64 go build -o ../../build/release/dwn-linux-amd64-$(SEMVER) .
	cd applications/dwn && GOOS=linux GOARCH=arm64 go build -o ../../build/release/dwn-linux-arm64-$(SEMVER) .
	cd applications/dwn && GOOS=linux GOARCH=386 go build -o ../../build/release/dwn-linux-386-$(SEMVER) .
	cd applications/dwn && GOOS=windows GOARCH=amd64 go build -o ../../build/release/dwn-windows-amd64-$(SEMVER) .
	cd applications/dwn && GOOS=windows GOARCH=386 go build -o ../../build/release/dwn-windows-386-$(SEMVER) .