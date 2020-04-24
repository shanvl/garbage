DOCKER_IMAGE_NAME=shanvl/garbage-rest-svc
GOOS?=darwin
GOARCH?=amd64

.PHONY: up
up:
	docker-compose -f ./docker/docker-compose.yml up -d --build

.PHONY: down
down:
	docker-compose -f ./docker/docker-compose.yml down

.PHONY: stop
stop:
	docker-compose -f ./docker/docker-compose.yml stop

.PHONY: test
test:
	docker-compose -f ./docker/docker-compose.test.yml up --build --abort-on-container-exit
	docker-compose -f ./docker/docker-compose.test.yml down --volumes

.PHONY: test-db-up
db-only-up:
	docker-compose -f ./docker/docker-compose.test.yml up -d --build db

.PHONY: test-db-down
db-only-down:
	docker-compose -f ./docker/docker-compose.test.yml down --volumes

.PHONY: build
build:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build \
		-o bin/restsvc ./cmd/restsvc

.PHONY: check
check: local-test vet

.PHONY: run
run: build
	./bin/restsvc

.PHONY: test
local-test:
	go test -v -race -timeout 30s ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: docker-build
docker-build:
	docker build -t ${DOCKER_IMAGE_NAME} -f ./docker/Dockerfile .

.PHONY: docker-push
docker-push:
	docker push ${DOCKER_IMAGE_NAME}

