DOCKER_IMAGE_NAME=shanvl/garbage-events-service
GOOS?=darwin
GOARCH?=amd64

build:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build \
		-o bin/eventsvc ./cmd/eventsvc

check: test vet

docker-build:
	docker build -t ${DOCKER_IMAGE_NAME} .

docker-push:
	docker push ${DOCKER_IMAGE_NAME}

docker-run:
	docker run --rm ${DOCKER_IMAGE_NAME}

run: build
	./bin/eventsvc

test:
	go test -v -race -timeout 30s ./...

vet:
	go vet ./...

.PHONY: build check docker-build docker-push docker-run run test vet
