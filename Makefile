DOCKER_IMAGE_NAME=shanvl/garbage-rest-svc
GOOS?=darwin
GOARCH?=amd64

up:
	docker-compose -f ./docker/docker-compose.yml up -d --build

down:
	docker-compose -f ./docker/docker-compose.yml down

stop:
	docker-compose -f ./docker/docker-compose.yml stop

test:
	docker-compose -f ./docker/docker-compose.test.yml up --build --abort-on-container-exit -V
	docker-compose -f ./docker/docker-compose.test.yml down --volumes

db-only-up:
	docker-compose -f ./docker/docker-compose.test.yml up -d --build db

db-only-down:
	docker-compose -f ./docker/docker-compose.test.yml down --volumes

build:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build \
		-o bin/eventssvc ./cmd/eventssvc

check: local-test vet

run: build
	./bin/eventssvc

local-test:
	go test -v -race -timeout 30s ./...

vet:
	go vet ./...

docker-build:
	docker build -t ${DOCKER_IMAGE_NAME} -f ./docker/Dockerfile .

docker-push:
	docker push ${DOCKER_IMAGE_NAME}

gen-events:
	 protoc --proto_path=api/events/v1/proto --proto_path=third_party --go_out=plugins=grpc:api/events/v1/pb \
	 --grpc-gateway_out=:api/events/v1/pb --openapiv2_out=:api/events/v1/swagger api/events/v1/proto/*.proto

gen-health:
	protoc --proto_path=api/health/v1/proto --go_out=plugins=grpc:api/health/v1/pb api/health/v1/proto/*.proto

gen-all: gen-events gen-health

.PHONY: up down stop test test-db-up test-db-down build check run test vet docker-build docker-push gen-events \
gen-health gen-all