GOOS?=darwin
GOARCH?=amd64

up:
	docker-compose -f ./docker/docker-compose.yml up -d

down:
	docker-compose -f ./docker/docker-compose.yml down

stop:
	docker-compose -f ./docker/docker-compose.yml stop

test:
	docker-compose -f ./docker/docker-compose.test.yml up --build --abort-on-container-exit -V
	docker-compose -f ./docker/docker-compose.test.yml down --volumes

db-only-down:
	docker-compose -f ./docker/docker-compose.test.yml down --volumes

build-authsvc:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build \
		-o bin/authsvc ./cmd/authsvc

build-eventsvc:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build \
		-o bin/eventsvc ./cmd/eventsvc

gen-all: gen-auth gen-events gen-health

gen-auth:
	 protoc --proto_path=api/auth/v1/proto --proto_path=third_party --go_out=plugins=grpc:api/auth/v1/pb \
	 --grpc-gateway_out=:api/auth/v1/pb --openapiv2_out=allow_merge=true:api/auth/v1/swagger api/auth/v1/proto/*.proto

gen-health:
	protoc --proto_path=api/health/v1/proto --go_out=plugins=grpc:api/health/v1/pb api/health/v1/proto/*.proto

gen-events:
	 protoc --proto_path=api/events/v1/proto --proto_path=third_party --go_out=plugins=grpc:api/events/v1/pb \
	 --grpc-gateway_out=:api/events/v1/pb --openapiv2_out=allow_merge=true:api/events/v1/swagger \
	 api/events/v1/proto/*.proto

cert:
	cd cert; ./gen.sh; cd ..

.PHONY: build-authsvc build-eventsvc cert down gen-all gen-auth gen-events gen-health stop test up