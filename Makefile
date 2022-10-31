APP_NAME=golang-grpc-base-project
APP_VERSION=latest
DOCKER_REGISTRY=ghcr.io/xdorro
DOCKER_IMAGE=$(DOCKER_REGISTRY)/$(APP_NAME):$(APP_VERSION)

docker.build:
	docker build -t $(DOCKER_IMAGE) .

docker.push:
	docker push $(DOCKER_IMAGE)

docker.dev: docker.build docker.push

docker.run:
	docker-compose -f docker-compose.yml up -d --force-recreate

wire.gen:
	wire ./...

buf.gen:
	buf generate

buf.update:
	cd proto/ && buf mod update

lint.run:
	golangci-lint run --fast ./...

go.install:
	go install github.com/google/wire/cmd/wire@v0.5.0

	go install github.com/bufbuild/buf/cmd/buf@v1.9.0
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.12.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.0

go.gen: wire.gen buf.gen

go.lint: lint.run

go.get:
	go get -u ./...

go.tidy:
	go mod tidy

go.test:
	go test ./...

jwt:
	openssl genrsa -out id_rsa 4096
	openssl rsa -in id_rsa -pubout -out id_rsa.pub

