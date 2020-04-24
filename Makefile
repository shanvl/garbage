DOCKER_IMAGE_NAME=shanvl/garbage-rest-svc
GOOS?=darwin
GOARCH?=amd64

.PHONY: build
build:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build \
		-o bin/restsvc ./cmd/restsvc

.PHONY: check
check: test vet

.PHONY: docker-build
docker-build:
	docker build -t ${DOCKER_IMAGE_NAME} -f ./docker/Dockerfile .

.PHONY: docker-push
docker-push:
	docker push ${DOCKER_IMAGE_NAME}

.PHONY: docker-run
docker-run:
	docker run --rm ${DOCKER_IMAGE_NAME}

.PHONY: run
run: build
	./bin/eventsvc

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: vet
vet:
	go vet ./...
