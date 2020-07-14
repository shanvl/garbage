FROM golang:1.13 as builder
RUN useradd -u 10001 notroot
WORKDIR /garbage
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 make build-authsvc

FROM alpine
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
USER notroot
COPY --from=builder /garbage/bin/authsvc /authsvc
COPY --from=builder /garbage/internal/authsvc/jwt/keys_test /keys/
CMD ["/authsvc"]
